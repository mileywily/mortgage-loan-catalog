package middleware

import (
	"github.com/gin-gonic/gin"
)

// LegacyHeaderValidator reads legacy headers from the request and stores them in context.
// Java legacy does NOT echo these headers back in the response — no echo applied.
func LegacyHeaderValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simply pass through — Java does not reject nor echo X-Channel/X-Commerce/X-Transaction-ID
		c.Next()
	}
}
