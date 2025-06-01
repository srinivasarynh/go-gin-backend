
package response

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
)

type APIResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   interface{} `json:"error,omitempty"`
}

func Success(c *gin.Context, status int, message string, data interface{}) {
    c.JSON(status, APIResponse{
        Success: true,
        Message: message,
        Data:    data,
    })
}

func Error(c *gin.Context, status int, message string) {
    c.JSON(status, APIResponse{
        Success: false,
        Message: message,
        Error:   message,
    })
}

func ValidationError(c *gin.Context, err error) {
    var errors []string
    
    if validationErrors, ok := err.(validator.ValidationErrors); ok {
        for _, fieldError := range validationErrors {
            switch fieldError.Tag() {
            case "required":
                errors = append(errors, fieldError.Field()+" is required")
            case "email":
                errors = append(errors, fieldError.Field()+" must be a valid email")
            case "min":
                errors = append(errors, fieldError.Field()+" must be at least "+fieldError.Param()+" characters")
            case "max":
                errors = append(errors, fieldError.Field()+" must be at most "+fieldError.Param()+" characters")
            default:
                errors = append(errors, fieldError.Field()+" is invalid")
            }
        }
    } else {
        errors = append(errors, err.Error())
    }

    c.JSON(http.StatusBadRequest, APIResponse{
        Success: false,
        Message: "Validation failed",
        Error:   strings.Join(errors, ", "),
    })
}
