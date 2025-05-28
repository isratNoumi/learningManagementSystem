package routes

import (
	"github.com/kataras/iris/v12"
	"learningManagementSystem/internal-LMS/controllers"
	"learningManagementSystem/internal-LMS/middlewares"
)

func SetupCourseRoutes(api iris.Party) {

	// Use the CORS middleware
	api.UseRouter(middlewares.AddCorsMiddleware())
	// Authentication Part
	authApi := api.Party("/v1")
	authApi.Post("/registration", controllers.CreateUserRecord)
	authApi.Post("/login", controllers.CheckAuthentication)
	authApi.Post("/reset-password", controllers.ResetPassword)

	protectedApi := api.Party("/api")
	protectedApi.Use(middlewares.AddCookieMiddleware(), middlewares.VerifyMiddleware())
	protectedApi.Get("/logout", controllers.Logout)

	courses1 := protectedApi.Party("/v1/courses")
	courses1.Get("/", controllers.FindAllCourses)
	courses1.Post("/", controllers.CreateCourse)
	courses1.Get("/{course_id:int}", controllers.FindAllCourseContentsById)
	courses1.Get("/{course_name:string}/modules", controllers.FindAllModules)
	courses1.Get("/{course_name:string}/modules/{module_name:string}/lessons", controllers.FindAllLessonsbyCourses)
	courses1.Get("/{course_name:string}/quizzes", controllers.FindUnitsdetailsByID)
	courses1.Get("/all-details", controllers.FindAllCourseDetails)
	courses1.Get("/suggestions", controllers.FindCourseDetailsforSuggestions)

	courses := api.Party("/v1/courses")
	courses.Post("/{course_id:int}/modules", controllers.AddNewModule)
	courses.Post("/modules/{module_id:int}/lessons", controllers.AddNewLesson)
	courses.Post("/lessons/{lesson_id:int}/units", controllers.AddNewUnit)
	courses.Delete("/{course_id:int}", controllers.DeleteCourses)
	courses.Delete("/{course_id:int}/modules/{module_id:int}", controllers.DeleteModules)
	courses.Delete("/modules/{module_id:int}/lessons/{lesson_id:int}", controllers.DeleteLessons)
	courses.Delete("/lessons/{lesson_id:int}/units/{unit_id:int}", controllers.DeleteUnits)
	courses.Put("/units/{unit_id:int}/videos", controllers.UpdateVideoUrl)
	courses.Put("/units/{unit_id:int}/contents", controllers.UpdateText)

	userApi := api.Party("/v1/users")
	userApi.Post("/{user_id:int}/responses", controllers.PostResponse)
	userApi.Post("/{user_id:int}/courses/{course_id:int}", controllers.EnrollToCourse)
	userApi.Post("/{user_id:int}/units/{unit_id:int}/views", controllers.PostviewTime)
	userApi.Get("/{user_id:int}/courses/{course_id:int}/progress", controllers.ViewPercentage)
	userApi.Get("/{user_id:int}/courses", controllers.FindEnrolledCourses)

	// Instructors Part
	instructorApi := protectedApi.Party("/v1/instructors")
	instructorApi.Get("/{instructor_id}/courses", controllers.FindCoursesByInstructorsName)

}
