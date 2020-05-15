package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"shorts/database"
	h "shorts/helper"
	"shorts/models"

	"github.com/gin-gonic/gin"
)

// GetShortlinksTop : Returns top 20 shortlinks
func GetShortlinksTop(c *gin.Context) {
	domainIndexes := make(map[string]int)

	var topDomains []models.TopDomainsResponseData
	var shortlinkUse models.ShortlinkUse

	if linksUses, err := shortlinkUse.UseCount(); err != nil {
		c.JSON(http.StatusBadRequest, h.NewResponseError(err))
	} else {
		for _, linkUse := range linksUses {
			parsedURL, err := url.Parse(linkUse.FullLink)
			if err != nil {
				continue
			}

			if _, exists := domainIndexes[parsedURL.Host]; !exists {
				domainIndexes[parsedURL.Host] = len(topDomains)
				topDomains = append(topDomains, models.TopDomainsResponseData{Website: parsedURL.Host, UsesCount: linkUse.UsesCount})
			} else {
				topDomains[domainIndexes[parsedURL.Host]].UsesCount += linkUse.UsesCount
			}
		}

		sort.Slice(topDomains, func(left, right int) bool {
			if topDomains[left].UsesCount == topDomains[right].UsesCount {
				return topDomains[left].Website > topDomains[right].Website
			}
			return topDomains[left].UsesCount > topDomains[right].UsesCount
		})

		upperLimit := 20
		if len(topDomains) < 20 {
			upperLimit = len(topDomains)
		}

		c.JSON(http.StatusOK, h.NewResponseOkWithData(topDomains[:upperLimit]))
	}
}

// GetShortlinksGraph : Retuns uses count groupped by day, hour and minute
func GetShortlinksGraph(c *gin.Context) {

	var result models.ShortlinksGraphResponseData = make(models.ShortlinksGraphResponseData)

	var shortlinksUse []models.ShortlinkUse
	if err := database.DB.Find(&shortlinksUse).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		fmt.Println(err)
	} else {
		for _, shortlinkUse := range shortlinksUse {
			models.AddUseToGraph(&result, shortlinkUse)
		}
		c.JSON(http.StatusOK, h.NewResponseOkWithData(result))
	}
}
