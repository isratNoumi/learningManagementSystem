package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"gorm.io/gorm"
	"learningManagementSystem/internal-LMS/database"
	models2 "learningManagementSystem/internal-LMS/models"
)

// FindCoursesByInstructorsName : Find Courses of a  specific instructor id
func FindCoursesByInstructorsName(c iris.Context) {
	claims, _ := jwt.Get(c).(*models2.Claims)
	if claims == nil || claims.Userid == 0 || claims.Username == "" {
		_ = c.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().
			Title("Unauthorized").
			Detail("Invalid or missing JWT claims"))
		return
	}
	if claims.Role != 2 {
		c.StatusCode(iris.StatusForbidden)
		c.JSON(iris.Map{"error": "You are not authorized to access this resource.Please login as a instructor."})
		return
	}
	standardClaims := jwt.GetVerifiedToken(c).StandardClaims
	expiresAtString := standardClaims.ExpiresAt().
		Format(c.Application().ConfigurationReadOnly().GetTimeFormat())
	timeLeft := standardClaims.Timeleft()
	c.Writef("user_id=%d\nusername=%s\nexpires at: %s\ntime left: %s\n",
		claims.Userid, claims.Username, expiresAtString, timeLeft)
	id, err := c.Params().GetInt("instructor_id")

	if err != nil {
		c.StatusCode(iris.StatusBadRequest)
		c.JSON(iris.Map{"error": "Instructor id is required"})
		return
	}
	var courses []models2.InstructorResponse

	err = database.DB.Debug().Table("courses c").Select("c.name as courses_name ,COUNT(p.users_id) as no_of_students,c.category as courses_category").
		Joins("inner join instructors i on c.id=i.courses_id").
		Joins("inner join users u on u.id=i.users_id").
		Joins("inner join progress_reports p on p.courses_id = c.id").Where("i.id=?", id).Group("c.id").Scan(&courses).Error

	//err = database.DB.Table("courses c").Joins("inner join instructors i on c.id=i.courses_id").
	//	Select("c.name as courses_name, c.category as courses_category").
	//	Where("i.id=?", id).Scan(&courses).Error
	if err != nil {
		c.StatusCode(iris.StatusInternalServerError)
		c.JSON(iris.Map{"error": "failed to query courses"})
	}

	c.StatusCode(iris.StatusOK)
	err = c.JSON(courses)
	if err != nil {
		c.JSON(iris.Map{
			"error": "Failed to serialize response: " + err.Error(),
		})
		return
	}

}

func CreateCourse(ctx iris.Context) {
	claims, _ := jwt.Get(ctx).(*models2.Claims)
	if claims == nil || claims.Userid == 0 || claims.Username == "" {
		_ = ctx.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().
			Title("Unauthorized").
			Detail("Invalid or missing JWT claims"))
		return
	}
	if claims.Role != 2 {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.JSON(iris.Map{"error": "You are not authorized to access this resource.Please login as a instructor."})
		return
	}

	var course models2.Course
	// Read the JSON body into the course struct
	if err := ctx.ReadJSON(&course); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid JSON format"})
		return
	}
	// Check if the course already exists
	var courseExists bool
	err := database.DB.Raw("SELECT EXISTS(SELECT 1 FROM courses WHERE name = ?)", course.Name).Scan(&courseExists).Error
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "failed to verify course: " + err.Error()})
		return

	}
	if courseExists {
		ctx.StatusCode(iris.StatusConflict)
		ctx.JSON(iris.Map{"error": "Course already exists"})
		return
	}
	txErr := database.DB.Transaction(func(tx *gorm.DB) error {

		tx = database.DB.Session(&gorm.Session{FullSaveAssociations: true})
		if err := tx.Create(&course).Error; err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(iris.Map{"error": "Failed to create course: " + err.Error()})
			return err
		}
		return nil
	})
	// Check if the transaction was successful
	if txErr != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "failed to create course: " + txErr.Error()})
	}
	// If the transaction is successful, return a success message

	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(iris.Map{"message": "Course created successfully"})
}

func DeleteCourses(ctx iris.Context) {
	id, err := ctx.Params().GetInt("course_id")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Course id is required"})
		return
	}
	var courseExist bool
	err = database.DB.Raw("SELECT EXISTS(SELECT 1 FROM courses WHERE id = ?)", id).Scan(&courseExist).Error
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "failed to verify course: " + err.Error()})
		return
	}
	if !courseExist {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{"error": "Course not found"})
		return
	}
	var course models2.Course
	err = database.DB.Where("id=?", id).Delete(&course).Error
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "failed to delete course"})
		return
	}
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{"message": "Course deleted successfully"})
}

func DeleteModules(ctx iris.Context) {
	courseId, err := ctx.Params().GetInt("course_id")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Course id is required"})
		return
	}
	id, err := ctx.Params().GetInt("module_id")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Module id is required"})
		return
	}
	var moduleExist bool
	err = database.DB.Raw("SELECT EXISTS(SELECT 1 FROM modules WHERE id = ? and courses_id=?)", id, courseId).Scan(&moduleExist).Error
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "failed to verify module: " + err.Error()})
		return
	}
	if !moduleExist {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{"error": "Module not found"})
		return
	}
	var module models2.Module
	err = database.DB.Where("id=? and courses_id=?", id, courseId).Delete(&module).Error
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "failed to delete module"})
		return
	}
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{"message": "Module deleted successfully"})
}
func DeleteLessons(ctx iris.Context) {

	moduleId, err := ctx.Params().GetInt("module_id")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Module id is required"})
		return
	}
	lessonId, err := ctx.Params().GetInt("lesson_id")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Lesson id is required"})
		return
	}
	var lessonExist bool
	err = database.DB.Raw("SELECT EXISTS(SELECT 1 FROM lessons WHERE id = ? and modules_id=?)", lessonId, moduleId).Scan(&lessonExist).Error
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "failed to verify lesson: " + err.Error()})
		return

	}
	if !lessonExist {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{"error": "Lesson not found"})
		return
	}
	var lesson models2.Lesson
	err = database.DB.Where("id=? and modules_id=?", lessonId, moduleId).Delete(&lesson).Error
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "failed to delete lesson"})
		return

	}
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{"message": "Lesson deleted successfully"})

}

func DeleteUnits(c iris.Context) {
	id, err := c.Params().GetInt("unit_id")
	if err != nil {
		c.StatusCode(iris.StatusBadRequest)
		c.JSON(iris.Map{"error": "Unit id is required"})
		return
	}
	lessonId, err := c.Params().GetInt("lesson_id")
	if err != nil {
		c.StatusCode(iris.StatusBadRequest)
		c.JSON(iris.Map{"error": "Lesson id is required"})
		return
	}
	var unitExist bool
	err = database.DB.Raw("SELECT EXISTS(SELECT 1 FROM units WHERE id = ? and lessons_id=?)", id, lessonId).Scan(&unitExist).Error
	if err != nil {
		c.StatusCode(iris.StatusInternalServerError)
		c.JSON(iris.Map{"error": "failed to verify unit: " + err.Error()})
		return

	}
	if !unitExist {
		c.StatusCode(iris.StatusNotFound)
		c.JSON(iris.Map{"error": "Unit not found"})
		return

	}
	var unit models2.Unit
	err = database.DB.Where("id=? and lessons_id=?", id, lessonId).Delete(&unit).Error
	if err != nil {
		c.StatusCode(iris.StatusInternalServerError)
		c.JSON(iris.Map{"error": "failed to delete unit"})
		return
	}
	c.StatusCode(iris.StatusOK)
	c.JSON(iris.Map{"message": "Unit deleted successfully"})

}

func AddNewModule(ctx iris.Context) {
	courseId, err := ctx.Params().GetInt("course_id")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Course id is required"})
		return
	}
	var courseExist bool
	err = database.DB.Raw("SELECT EXISTS(SELECT 1 FROM courses WHERE id = ?)", courseId).Scan(&courseExist).Error
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "failed to verify course: " + err.Error()})
		return
	}
	if !courseExist {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{"error": "Course not found"})
		return
	}

	// Read the JSON body into the module struct
	var module models2.Module
	if err := ctx.ReadJSON(&module); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid JSON format"})
		return
	}
	module.CoursesID = int64(courseId)
	txErr := database.DB.Transaction(func(tx *gorm.DB) error {

		tx = database.DB.Session(&gorm.Session{FullSaveAssociations: true})
		if err := tx.Create(&module).Error; err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(iris.Map{"error": "Failed to create course: " + err.Error()})
			return err
		}
		return nil
	})
	// Check if the transaction was successful
	if txErr != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "failed to create module: " + txErr.Error()})
	}
	// If the transaction is successful, return a success message

	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(iris.Map{"message": "Module created successfully"})

}
func AddNewLesson(ctx iris.Context) {
	moduleId, err := ctx.Params().GetInt("module_id")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Module id is required"})
		return
	}
	var moduleExist bool
	err = database.DB.Raw("SELECT EXISTS(SELECT 1 FROM modules WHERE id = ?)", moduleId).Scan(&moduleExist).Error
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "failed to verify module: " + err.Error()})
		return
	}
	if !moduleExist {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{"error": "Module not found"})
		return
	}

	var lesson models2.Lesson
	if err := ctx.ReadJSON(&lesson); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid JSON format"})
		return
	}
	lesson.ModulesID = int64(moduleId)
	txErr := database.DB.Transaction(func(tx *gorm.DB) error {

		tx = database.DB.Session(&gorm.Session{FullSaveAssociations: true})
		if err := tx.Create(&lesson).Error; err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(iris.Map{"error": "Failed to create lesson: " + err.Error()})
			return err
		}
		return nil
	})
	if txErr != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "failed to create lesson: " + txErr.Error()})
	}
	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(iris.Map{"message": "lesson created successfully"})
}
func AddNewUnit(ctx iris.Context) {
	lessonId, err := ctx.Params().GetInt("lesson_id")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Lesson id is required"})
		return
	}
	var lessonExist bool
	err = database.DB.Raw("SELECT EXISTS(SELECT 1 FROM lessons WHERE id = ?)", lessonId).Scan(&lessonExist).Error
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "failed to verify lesson: " + err.Error()})
		return
	}
	if !lessonExist {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{"error": "Lesson not found"})
		return
	}

	var unit models2.Unit
	if err := ctx.ReadJSON(&unit); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid JSON format"})
		return
	}
	unit.LessonsID = int64(lessonId)
	txErr := database.DB.Transaction(func(tx *gorm.DB) error {

		tx = database.DB.Session(&gorm.Session{FullSaveAssociations: true})
		if err := tx.Create(&unit).Error; err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(iris.Map{"error": "Failed to create unit: " + err.Error()})
			return err
		}
		return nil
	})
	if txErr != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "failed to create unit: " + txErr.Error()})
	}
	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(iris.Map{"message": "Unit created successfully"})
}

// UpdateVideoUrl updates the video URL for a specific unit
func UpdateVideoUrl(ctx iris.Context) {
	unitId, err := ctx.Params().GetInt("unit_id")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Unit id is required"})
		return
	}
	var unitExist bool
	err = database.DB.Raw("SELECT EXISTS(SELECT 1 FROM units WHERE id = ?)", unitId).Scan(&unitExist).Error
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "failed to verify unit: " + err.Error()})
		return

	}
	if !unitExist {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{"error": "Unit not found"})
		return
	}

	var unitdetails models2.Video
	if err := ctx.ReadJSON(&unitdetails); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid JSON format"})
		return
	}
	txErr := database.DB.Transaction(func(tx *gorm.DB) error {
		// Update course fields

		var detailsID int
		var detailsID2 int
		err := tx.Table("units_details").
			Joins("INNER JOIN units_fields ON units_fields.id = units_details.units_fields_id").
			Joins("INNER JOIN units ON units_fields.units_id = units.id").
			Where("units.id = ? AND units_fields.fields = ?", unitId, "video_url").
			Pluck("units_details.id", &detailsID).Error

		if err != nil {
			return err
		}
		//  Update the description for the found units_details ID(s)
		err = tx.Model(&models2.UnitsDetails{}).
			Where("id = ?", detailsID).
			Update("description", unitdetails.URL).Error
		if err != nil {
			return err
		}
		err = tx.Table("units_details").
			Joins("INNER JOIN units_fields ON units_fields.id = units_details.units_fields_id").
			Joins("INNER JOIN units ON units_fields.units_id = units.id").
			Where("units.id = ? AND units_fields.fields = ?", unitId, "duration").
			Pluck("units_details.id", &detailsID2).Error

		if err != nil {
			return err
		}

		//  Update the duration for the found units_details ID(s)
		err = tx.Model(&models2.UnitsDetails{}).
			Where("id = ?", detailsID2).
			Update("description", unitdetails.Duration).Error
		if err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Update failed: " + txErr.Error()})
		return
	}
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{"message": "Video URL updated successfully"})

}
func UpdateText(ctx iris.Context) {
	unitId, err := ctx.Params().GetInt("unit_id")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Unit id is required"})
		return
	}
	var unitExist bool
	err = database.DB.Raw("SELECT EXISTS(SELECT 1 FROM units WHERE id = ?)", unitId).Scan(&unitExist).Error
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "failed to verify unit: " + err.Error()})
		return

	}
	if !unitExist {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{"error": "Unit not found"})
		return
	}
	var unitdetails models2.Content
	if err := ctx.ReadJSON(&unitdetails); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid JSON format"})
		return
	}
	txErr := database.DB.Transaction(func(tx *gorm.DB) error {
		var detailsID int
		var detailsID2 int
		err := tx.Table("units_details").
			Joins("INNER JOIN units_fields ON units_fields.id = units_details.units_fields_id").
			Joins("INNER JOIN units ON units_fields.units_id = units.id").
			Where("units.id = ? AND units_fields.fields = ?", unitId, "content").
			Pluck("units_details.id", &detailsID).Error

		if err != nil {
			return err
		}
		//  Update the description for the found units_details ID(s)
		err = tx.Model(&models2.UnitsDetails{}).
			Where("id = ?", detailsID).
			Update("description", unitdetails.Text).Error
		if err != nil {
			return err
		}
		err = tx.Table("units_details").
			Joins("INNER JOIN units_fields ON units_fields.id = units_details.units_fields_id").
			Joins("INNER JOIN units ON units_fields.units_id = units.id").
			Where("units.id = ? AND units_fields.fields = ?", unitId, "length").
			Pluck("units_details.id", &detailsID2).Error
		if err != nil {
			return err

		}
		//  Update the length for the found units_details ID(s)
		err = tx.Model(&models2.UnitsDetails{}).
			Where("id = ?", detailsID2).
			Update("description", unitdetails.Length).Error
		if err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Update failed: " + txErr.Error()})
		return
	}
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{"message": "Content updated successfully"})

}

//func UpdateCourseWithStructure(ctx iris.Context) {
//	courseID, _ := ctx.Params().GetInt("course_id")
//	var input models2.Course
//	if err := ctx.ReadJSON(&input); err != nil {
//		ctx.StatusCode(iris.StatusBadRequest)
//		ctx.JSON(iris.Map{"error": "Invalid JSON"})
//		return
//	}
//
//	txErr := database.DB.Transaction(func(tx *gorm.DB) error {
//		// Update course fields
//		if err := tx.Model(&models2.Course{}).Where("id = ?", courseID).Updates(input).Error; err != nil {
//			return err
//		}
//		// Update or create modules
//		for _, m := range input.Modules {
//			m.CoursesID = int64(courseID)
//			if m.ID == 0 {
//				if err := tx.Create(&m).Error; err != nil {
//					return err
//				}
//			} else {
//				if err := tx.Model(&models2.Module{}).Where("id = ?", m.ID).Updates(m).Error; err != nil {
//					return err
//				}
//			}
//			// Update or create lessons for each module
//			for _, l := range m.Lessons {
//				l.ModulesID = m.ID
//				if l.ID == 0 {
//					if err := tx.Create(&l).Error; err != nil {
//						return err
//					}
//				} else {
//					if err := tx.Model(&models2.Lesson{}).Where("id = ?", l.ID).Updates(l).Error; err != nil {
//						return err
//					}
//				}
//				for _, u := range l.Units {
//					u.LessonsID = l.ID
//					if u.ID == 0 {
//						if err := tx.Create(&u).Error; err != nil {
//							return err
//						}
//					} else {
//						if err := tx.Model(&models2.Unit{}).Where("id = ?", u.ID).Updates(u).Error; err != nil {
//							return err
//						}
//					}
//					for _, uf := range u.UnitsFields {
//						uf.UnitsID = u.ID
//						if uf.ID == 0 {
//							if err := tx.Create(&uf).Error; err != nil {
//								return err
//							}
//						} else {
//							if err := tx.Model(&models2.UnitsFields{}).Where("id = ?", uf.ID).Updates(uf).Error; err != nil {
//								return err
//							}
//						}
//					}
//				}
//
//			}
//		}
//
//		return nil
//	})
//
//	if txErr != nil {
//		ctx.StatusCode(iris.StatusInternalServerError)
//		ctx.JSON(iris.Map{"error": "Update failed: " + txErr.Error()})
//		return
//	}
//	ctx.StatusCode(iris.StatusOK)
//	ctx.JSON(iris.Map{"message": "Course and structure updated"})
//}
