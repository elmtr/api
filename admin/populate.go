package admin

import (
	"api/grip"
	"api/utils"
	"strconv"
	"strings"

	"github.com/deta/deta-go/service/base"
	"github.com/gofiber/fiber/v2"
)

func populate(g fiber.Router) {
  admin := g.Group("/populate")

  admin.Get("/test", func (c *fiber.Ctx) error {
    return c.SendString("Hello world")
  })

  admin.Post("/grades", authMiddleware, func (c *fiber.Ctx) error {
    var admin grip.Admin
    utils.GetLocals(c, "admin", &admin)

    var grades []grip.Grade
    for gradeNumber, gradeLetters := range subjectsHighSchool {
      for gradeLetter := range gradeLetters {
        gradeNumberInt, _ := strconv.Atoi(gradeNumber) 

        grade := grip.Grade {
          GradeNumber: gradeNumberInt,
          GradeLetter: gradeLetter,
          SchoolKey: admin.SchoolKey,
          YearFrom: 2022,
          YearTo: 2023,
        }

        // grade.Put()
        grades = append(grades, grade)
      }
    }

    return c.JSON(grades)
  })

  admin.Post("/subjects", authMiddleware, func (c *fiber.Ctx) error {
    var subjectsObjects []grip.Subject
    for gradeNumber, gradeLetters := range subjectsHighSchool {
      for gradeLetter, subjectsString := range gradeLetters {
        gradeNumberInt, _ := strconv.Atoi(gradeNumber) 

        grade, _ := grip.GetGrade(base.Query {
          {
            "gradeNumber": gradeNumberInt,
            "gradeLetter": gradeLetter,
          },
        })
        subjects := strings.Split(subjectsString, ", ")
        for ord, subject := range subjects {
          subjectID := utils.GenID()
          subjectObject := grip.Subject {
            Name: subject,
            Key: subjectID,
            Grade: grade,
            Ord: ord + 1,
          }
          // subjectObject.Put()
          subjectsObjects = append(subjectsObjects, subjectObject)
        }
      }
    }
    return c.JSON(subjectsObjects)
  })
}