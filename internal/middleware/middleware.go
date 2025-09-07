package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/onunkwor/flypro-assestment-v2/internal/repository"
	"github.com/onunkwor/flypro-assestment-v2/internal/utils"
)

func ReportOwnershipMiddleware(reportRepo repository.ReportRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Query("userID")
		reportIDStr := c.Param("id")
		userID, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil || userID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
			c.Abort()
			return
		}

		reportID, err := strconv.ParseUint(reportIDStr, 10, 64)
		if err != nil || reportID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid report ID"})
			c.Abort()
			return
		}

		report, err := reportRepo.GetExpenseReportByID(c.Request.Context(), uint(reportID))
		if err != nil {
			utils.BadRequestResponse(c, "failed to retrieve report")
			c.Abort()
			return
		}

		if report.UserID != uint(userID) {
			utils.ForbiddenResponse(c, "you do not have permission to access this resource")
			c.Abort()
			return
		}

		c.Set("userID", uint(userID))
		c.Set("reportID", uint(reportID))
		c.Next()
	}
}
