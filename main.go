package main

import (
	"encoding/json"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/x/errors"
	"gorm.io/gorm"
	"learningManagementSystem/models"
	"log"
	"net/http"
	_ "net/http"
	"net/url"
	"strconv"
	"strings"
)

func main() {
	app := iris.New()
	err := models.InitDatabase()
	if err != nil {
		// Log the error and exit
		log.Fatalln("could not create database", err)
	}
	/// users part
	app.Get("/v1/courses", FindAllCourses)
	app.Get("/v1/courses/{course_name}/modules", FindAllModules)
	app.Get("/v1/courses/{course_name}/{module_name}/lessons", FindAllLessonsbyCourses)
	app.Get("/v1/courses/{id}", FindAllCourseContentsById)
	app.Get("/v1/courses/details", FindAllCourseDetails)
	app.Get("/v1/courses/suggestions", FindCourseDetailsforSuggestions)
	app.Get("/v1/courses/{id}/noofstudents", FindEnrolledStudents)
	app.Get("/v1/courses/{name}/quizes", FindUnitsdetailsByID)

	// Instructors Part
	app.Get("/v1/instructors/{instructor_id}/courses", FindCoursesByInstructorsName)
	app.Listen(":8086")

}

// FindAllCourses :show all the courses : pagination added : sorting( name , created_at, total_score ) added : search by name , category implemented
func FindAllCourses(c iris.Context) {
	pageStr := c.URLParamDefault("page", "0")
	sizeStr := c.URLParamDefault("size", "10")
	sortStr := c.URLParamDefault("sort", "created_at:desc") // Default sort
	filterStr := c.URLParam("filter")

	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil || page < 0 {
		c.JSON(http.StatusBadRequest)
		c.JSON(iris.Map{"error": "Invalid page parameter"})
		return
	}

	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil || size <= 0 {
		c.JSON(http.StatusBadRequest)
		c.JSON(iris.Map{"error": "Invalid size parameter"})
		return
	}

	// Parse and validate sort
	sortParts := strings.Split(sortStr, ":")
	if len(sortParts) != 2 {
		c.JSON(http.StatusBadRequest)
		c.JSON(iris.Map{"error": "Invalid sort format. Use field:order (e.g., name:asc)"})
		return
	}
	sortField, sortOrder := sortParts[0], sortParts[1]

	// Validate sort field
	allowedFields := map[string]string{
		"name":        "name",
		"total_score": "total_score",
		"created_at":  "created_at",
	}
	dbField, ok := allowedFields[sortField]
	if !ok {
		c.JSON(http.StatusBadRequest)
		c.JSON(iris.Map{"error": "Invalid sort field. Allowed: name, total_score, created_at"})
		return
	}

	// Validate sort order
	if sortOrder != "asc" && sortOrder != "desc" {
		c.JSON(http.StatusBadRequest)
		c.JSON(iris.Map{"error": "Invalid sort order. Use asc or desc"})
		return
	}

	// Parse and validate filter
	var filters [][]string
	if filterStr != "" {
		if err := json.Unmarshal([]byte(filterStr), &filters); err != nil {
			err := c.JSON(http.StatusBadRequest)
			if err != nil {
				c.StatusCode(http.StatusInternalServerError)
				c.WriteString("Failed to serialize error response: " + err.Error())

			}
			c.JSON(iris.Map{"error": "Invalid filter format. Use JSON array of [field,value] pairs (e.g., [[\"category\",\"Programming\"]])"})
			return
		}
		for _, f := range filters {
			if len(f) != 2 {
				err := c.JSON(http.StatusBadRequest)
				if err != nil {
					c.StatusCode(http.StatusInternalServerError)
					c.WriteString("Failed to serialize error response: " + err.Error())

				}
				c.JSON(iris.Map{"error": "Invalid filter entry. Each entry must be [field,value]"})
				return
			}
			field := f[0]
			if field != "category" && field != "name" {
				err := c.JSON(http.StatusBadRequest)
				if err != nil {
					c.StatusCode(http.StatusInternalServerError)
					c.WriteString("Failed to serialize error response: " + err.Error())
				}
				c.JSON(iris.Map{"error": "Invalid filter field. Allowed: category, name"})
				return
			}
		}
	}

	// Calculate offset
	offset := page * size

	// Build query with filters
	query := models.DB.Model(&models.Course{})
	for _, f := range filters {
		field, value := f[0], f[1]
		switch field {
		case "category":
			query = query.Where("category = ?", value)
		case "name":
			query = query.Where("name LIKE ?", "%"+value+"%")
		}
	}
	// Get total count models.DB.Model(&models.Course{})
	var totalRecords int64
	err = query.Count(&totalRecords).Error
	if err != nil {
		log.Println("Error counting records:", err)
		c.JSON(http.StatusInternalServerError)
		c.JSON(iris.Map{"error": "Internal server error"})
		return
	}

	var course []models.Course
	err = query.Limit(int(size)).Offset(int(offset)).Order(dbField + " " + strings.ToUpper(sortOrder)).Find(&course).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("Error querying courses:", err.Error())
		c.JSON(iris.Map{"error": "Failed to fetch courses" + err.Error()})
		return
	}
	// Handle empty results
	if len(course) == 0 && page > 0 {
		if err = c.JSON(http.StatusNotFound); err != nil {
			c.StatusCode(http.StatusInternalServerError)
			c.WriteString("Failed to serialize error response: " + err.Error())
			return
		}
		c.JSON(iris.Map{"error": "No courses found for this page"})

		return
	}

	// Calculate total pages
	totalPages := (totalRecords + size - 1) / size
	courseDTOs := make([]models.CourseDTO, len(course))
	for i, course := range course {
		courseDTOs[i] = models.CourseDTO{
			ID:         course.ID,
			Name:       course.Name,
			Category:   course.Category,
			TotalScore: course.TotalScore,
			CreatedAt:  course.CreatedAt,
			UpdatedAt:  course.UpdatedAt,
		}
	}

	filterQuery := ""
	if filterStr != "" {
		filterQuery = "&filter=" + url.QueryEscape(filterStr)
	}
	sortQuery := "&sort=" + url.QueryEscape(sortStr)
	// Build HATEOAS links for the collection
	links := []models.Link{
		{Rel: "self", Href: "/v1/courses?page=" + strconv.FormatInt(page, 10) + "&size=" + strconv.FormatInt(size, 10) + sortQuery + filterQuery},
	}
	if page > 0 {
		links = append(links, models.Link{Rel: "prev", Href: "/v1/courses?page=" + strconv.FormatInt(page-1, 10) + "&size=" + strconv.FormatInt(size, 10) + sortQuery + filterQuery})
	}
	if page < totalPages-1 {
		links = append(links, models.Link{Rel: "next", Href: "/v1/courses?page=" + strconv.FormatInt(page+1, 10) + "&size=" + strconv.FormatInt(size, 10) + sortQuery + filterQuery})
	}
	links = append(links,
		models.Link{Rel: "first", Href: "/v1/courses?page=0&size=" + strconv.FormatInt(size, 10) + sortQuery + filterQuery},
		models.Link{Rel: "last", Href: "/v1/courses?page=" + strconv.FormatInt(totalPages-1, 10) + "&size=" + strconv.FormatInt(size, 10) + sortQuery + filterQuery},
	)

	// Build response
	response := models.Response{
		Data: courseDTOs,
		Pagination: models.Pagination{
			CurrentPage:  int(page),
			PageSize:     int(size),
			TotalRecords: int(totalRecords),
			TotalPages:   int(totalPages),
		},
		Links: links,
	}

	err = c.JSON(response)
	if err != nil {
		c.JSON(iris.Map{
			"error": "Failed to serialize response: " + err.Error(),
		})
		return
	}
	c.StatusCode(iris.StatusOK)
}

// FindAllModules :Find All the modules details by a course name
func FindAllModules(c iris.Context) {
	courseName := c.Params().Get("course_name")
	//c.JSON(courseName)
	var module []models.Module
	if courseName == "" {
		c.StatusCode(iris.StatusBadRequest)
		c.JSON(iris.Map{"error": "Course name is required"})
		return
	}
	err := models.DB.Debug().Table("modules m").Joins("inner join courses c on c.id=m.courses_id ").Where("c.name = ?", courseName).Find(&module).Error
	if err != nil {
		log.Println("Error querying modules:", err.Error())
		c.JSON(iris.Map{"error": "Failed to fetch modules" + err.Error()})
		return
	}
	c.StatusCode(iris.StatusOK)
	err = c.JSON(module)
	if err != nil {
		c.JSON(iris.Map{
			"error": "Failed to serialize response: " + err.Error(),
		})
		return
	}

}

// FindAllLessonsbyCourses :Find All the Lesson details by a course name and module name
func FindAllLessonsbyCourses(c iris.Context) {
	courseName := c.Params().Get("course_name")
	moduleName := c.Params().Get("module_name")

	if courseName == "" {
		c.StatusCode(iris.StatusBadRequest)
		c.JSON(iris.Map{"error": "Course name is required"})
		return
	}
	if moduleName == "" {
		c.StatusCode(iris.StatusBadRequest)
		c.JSON(iris.Map{"error": "Module name is required"})
		return
	}
	var lesson []models.Lesson
	err := models.DB.Table("lessons l").Joins("inner join modules m on l.modules_id=m.id").Joins("inner join courses c on m.courses_id=c.id ").
		Where("c.name = ?", courseName).Where("m.name = ?", moduleName).Find(&lesson).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("Error querying lessons:", err.Error())
		c.JSON(iris.Map{"error": "Failed to fetch lessons" + err.Error()})
		return
	}
	c.StatusCode(iris.StatusOK)
	err = c.JSON(lesson)
	if err != nil {
		c.JSON(iris.Map{
			"error": "Failed to serialize response: " + err.Error(),
		})
		return
	}

}

// FindEnrolledStudents :Find number of Enrolled Students in a specific course
func FindEnrolledStudents(c iris.Context) {
	var noofstudents int64
	id := c.Params().Get("id")
	if id == "" {
		c.StatusCode(iris.StatusBadRequest)
		c.JSON(iris.Map{"error": "Course id is required"})
		return
	}

	err1 := models.DB.Table("progress_report p").Select(" COUNT(p.id) as enrolled_students").
		Joins("INNER JOIN users u on p.user_id=u.id").
		Where("p.course_id = ?", id).Scan(&noofstudents).Error
	if err1 != nil {

		c.JSON(http.StatusBadRequest)
		c.JSON(iris.Map{"error": err1.Error()})

	}
	c.StatusCode(iris.StatusOK)
	err := c.JSON(iris.Map{"enrolled_students": noofstudents})
	if err != nil {
		c.JSON(iris.Map{
			"error": "Failed to serialize response: " + err.Error(),
		})
		return
	}

}

// FindAllCourseContentsById :Find ALl Course Contents by a course ID
func FindAllCourseContentsById(c iris.Context) {
	id := c.Params().Get("id")
	if id == "" {
		c.StatusCode(iris.StatusBadRequest)
		c.JSON(iris.Map{"error": "Course id is required"})
		return
	}
	// Query to fetch courses with modules and lessons
	var courses []models.Course
	err := models.DB.Preload("Modules", func(db *gorm.DB) *gorm.DB { return db.Select("id, name,courses_id") }).
		Preload("Modules.Lessons", func(db *gorm.DB) *gorm.DB { return db.Select("id, name,modules_id") }).
		Preload("Modules.Lessons.Units", func(db *gorm.DB) *gorm.DB { return db.Select("id,lessons_id,type") }).
		Preload("Modules.Lessons.Units.UnitsFields", func(db *gorm.DB) *gorm.DB { return db.Select("id, units_id,units_fields") }).
		Preload("Modules.Lessons.Units.UnitsFields.UnitsDetails", func(db *gorm.DB) *gorm.DB {
			return db.Select("units_details.id, units_details.units_id, units_details.units_fields_id, units_details.description")

		}).
		Preload("Modules.Lessons.Units.UnitsFields.UnitsDetails.QuizType", func(db *gorm.DB) *gorm.DB {
			return db.Joins(" inner JOIN units_details  ON units_details.id = quiz_types.units_details_id").
				Joins("JOIN units ON units.id = units_details.units_id").
				Where("units.type = ?", "quiz").
				Select("quiz_types.id, quiz_types.units_details_id, quiz_types.type")
		}).
		Preload("Modules.Lessons.Units.UnitsFields.UnitsDetails.QuizOption", func(db *gorm.DB) *gorm.DB {
			return db.Joins(" inner JOIN units_details  ON units_details.id = quiz_options.units_details_id").
				Joins("JOIN units ON units.id = units_details.units_id").
				Where("units.type = ?", "quiz").
				Select("quiz_options.id, quiz_options.units_details_id, quiz_options.options")
		}).
		Where("id=?", id).Find(&courses).Error

	if err != nil {

		c.StatusCode(iris.StatusInternalServerError)
		c.JSON(iris.Map{"error": "failed to query courses"})
		return
	}
	c.StatusCode(iris.StatusOK)
	// Return JSON response
	err = c.JSON(courses)
	if err != nil {
		c.JSON(iris.Map{
			"error": "Failed to serialize response: " + err.Error(),
		})
		return
	}

}

// FindAllCourseDetails :Find All Course Details
func FindAllCourseDetails(c iris.Context) {
	var instructor []models.InstructorResponse1
	err1 := models.DB.Table("courses c").Select(" c.name as course_name ,COUNT(p.user_id) as no_of_students ,u.name as instructor_name").
		Joins("left join instructors i on c.id=i.courses_id").
		Joins("left join users u on u.id=i.user_id").
		Joins("left join progress_report p on p.course_id = c.id").Group("c.id").Scan(&instructor).Error

	if err1 != nil {

		c.JSON(http.StatusBadRequest)
		c.JSON(iris.Map{"error": err1.Error()})

	}
	c.StatusCode(iris.StatusOK)
	err := c.JSON(instructor)
	if err != nil {
		c.JSON(iris.Map{
			"error": "Failed to serialize response: " + err.Error(),
		})
		return
	}

}

// FindCourseDetailsforSuggestions :Find Course Details for Suggestions
func FindCourseDetailsforSuggestions(c iris.Context) {
	var instructor []models.InstructorResponse1
	err1 := models.DB.Table("courses c").Select(" c.name as course_name ,COUNT(p.user_id) as no_of_students ,u.name as instructor_name").
		Joins("left join instructors i on c.id=i.courses_id").
		Joins("left join users u on u.id=i.user_id").
		Joins("left join progress_report p on p.course_id = c.id").Group("c.id").Order("rand()").Limit(5).Scan(&instructor).Error

	if err1 != nil {

		c.JSON(http.StatusBadRequest)
		c.JSON(iris.Map{"error": err1.Error()})

	}
	c.StatusCode(iris.StatusOK)
	err := c.JSON(instructor)
	if err != nil {
		c.JSON(iris.Map{
			"error": "Failed to serialize response: " + err.Error(),
		})
		return
	}

}

// FindCoursesByInstructorsName : Find Courses of a instructor id
func FindCoursesByInstructorsName(c iris.Context) {
	id := c.Params().Get("instructor_id")
	var courses []models.InstructorResponse
	if id == "" {
		c.StatusCode(iris.StatusBadRequest)
		c.JSON(iris.Map{"error": "Instructor id is required"})
		return
	}
	err := models.DB.Table("courses c").Joins("inner join instructors i on c.id=i.courses_id").
		Select("c.name as courses_name, c.category as courses_category").
		Where("i.id=?", id).Scan(&courses).Error
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

//  SELECT  uf.units_fields as Content_Type, ud.description FROM  courses c inner JOIN    modules m ON c.id = m.courses_id inner JOIN
// lessons l ON m.id = l.modules_id inner JOIN    units u ON l.id = u.lessons_id inner JOIN  units_fields uf ON u.id = uf.units_id
// inner JOIN units_details ud ON uf.id = ud.units_fields_id inner JOIN  quize_types q ON ud.id = q.units_details_id  WHERE    c.id = 5 and q.type = "multiple_choice"

func FindUnitsdetailsByID(c iris.Context) {
	name := c.Params().Get("name")
	typename := "multiple_choice"
	typename2 := "true_false"
	var units []models.MCQuiz
	err1 := models.DB.Table("courses c").Select("l.name as lesson_name,u.type as unit_type,q.type as question_type, uf.units_fields as content_type, ud.description as question").
		Joins("INNER JOIN modules m ON c.id = m.courses_id").Joins(" INNER JOIN lessons l ON m.id = l.modules_id").
		Joins("inner JOIN units u ON l.id = u.lessons_id ").Joins("inner JOIN  units_fields uf ON u.id = uf.units_id").
		Joins("inner JOIN units_details ud ON uf.id = ud.units_fields_id ").Joins("inner JOIN  quize_types q ON ud.id = q.units_details_id").
		Where("q.type =? OR q.type=?", typename, typename2).
		Where("c.name = ?", name).Scan(&units).Error
	if err1 != nil {

		c.JSON(http.StatusBadRequest)
		c.JSON(iris.Map{"error": err1.Error()})

	}
	err := c.JSON(units)
	if err != nil {
		c.JSON(iris.Map{
			"error": "Failed to serialize response: " + err.Error(),
		})
		return
	}
	c.StatusCode(iris.StatusOK)

}
