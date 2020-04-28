package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"shorts/database"
	h "shorts/helper"
	"shorts/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetShortlinks : Send all short links of current user
func GetShortlinks(c *gin.Context) {
	userID := c.MustGet(gin.AuthUserKey).(uint64)

	var shortlinks []models.Shortlink
	var shortlinksResponse []models.ShortlinkResponseData
	if err := database.DB.Preload("Uses").Where(&models.Shortlink{OwnerID: userID}).Find(&shortlinks).Error; err != nil {
		c.JSON(http.StatusBadRequest, h.NewResponseError(err))
	} else {
		for _, item := range shortlinks {
			shortlinksResponse = append(shortlinksResponse, models.ShortlinkResponseData{
				ID:    item.ID,
				Short: item.Short,
				Full:  item.Full,
			})
		}
		c.JSON(http.StatusOK, h.NewResponseOkWithData(shortlinksResponse))
	}
}

// AddShortlink : Create short link
func AddShortlink(c *gin.Context) {
	userID := c.MustGet(gin.AuthUserKey).(uint64)

	var shortlinkData models.ShortlinkAddData

	if err := c.ShouldBindJSON(&shortlinkData); err != nil {
		c.JSON(http.StatusBadRequest, h.NewResponseError(err))
		return
	}

	if parsedURL, err := url.Parse(shortlinkData.Full); err != nil {
		c.JSON(http.StatusBadRequest, h.NewResponseError(err))
	} else {
		if !parsedURL.IsAbs() {
			c.JSON(http.StatusBadRequest, h.NewResponseError(h.NewAbsoluteLinksOnlyError()))
			return
		}

		shortlink := models.Shortlink{
			OwnerID: userID,
			Short:   shortlinkData.Short,
			Full:    shortlinkData.Full,
		}

		if dbc := database.DB.Create(&shortlink); dbc.Error != nil {
			c.JSON(http.StatusBadRequest, h.NewResponseError(dbc.Error))
		} else {
			c.JSON(http.StatusCreated, h.NewResponseOkWithData(models.ShortlinkResponseData{
				ID:    shortlink.ID,
				Full:  shortlink.Full,
				Short: shortlink.Short,
			}))
		}
	}
}

// DeleteShortlink : Delete short link with the specified ID
func DeleteShortlink(c *gin.Context) {
	userID := c.MustGet(gin.AuthUserKey).(uint64)

	var shortlink models.Shortlink

	if parseResult, err := strconv.ParseUint(c.Params.ByName("id"), 10, 64); err != nil {
		c.JSON(http.StatusBadRequest, h.NewResponseError(err))
	} else {
		shortlinkID := parseResult

		if err := database.DB.Where(&models.Shortlink{ID: shortlinkID, OwnerID: userID}).First(&shortlink).Error; err != nil {
			c.JSON(http.StatusNotFound, h.NewResponseError(err))
		} else {
			if dbc := database.DB.Delete(&shortlink); dbc.Error != nil {
				c.JSON(http.StatusBadRequest, h.NewResponseError(err))
			} else {
				c.JSON(http.StatusOK, h.NewResponseOK())
			}
		}
	}
}

// GetShortlinkInfo : Send information about short link with the specified ID (including uses)
func GetShortlinkInfo(c *gin.Context) {
	userID := c.MustGet(gin.AuthUserKey).(uint64)

	var shortlink models.Shortlink
	if parseResult, err := strconv.ParseUint(c.Params.ByName("id"), 10, 64); err != nil {
		c.JSON(http.StatusBadRequest, h.NewResponseError(err))
	} else {
		shortlinkID := parseResult

		if err := database.DB.Preload("Uses").Where(&models.Shortlink{ID: shortlinkID, OwnerID: userID}).Find(&shortlink).Error; err != nil {
			c.JSON(http.StatusNotFound, h.NewResponseError(err))
		} else {
			c.JSON(http.StatusOK, h.NewResponseOkWithData(shortlink))
		}
	}
}

// GetShortlinkRedirect : Redirects to a full link by a short link
func GetShortlinkRedirect(c *gin.Context) {
	var shortlink models.Shortlink

	if shortlinkID, err := strconv.ParseUint(c.Params.ByName("short"), 36, 64); err != nil {
		c.JSON(http.StatusBadRequest, h.NewResponseError(err))
	} else {
		if err := database.DB.Preload("Uses").First(&shortlink, shortlinkID).Error; err != nil {
			c.JSON(http.StatusBadRequest, h.NewResponseError(err))
		} else {
			shortlink.Uses = append(shortlink.Uses, models.ShortlinkUse{LinkID: shortlink.ID, UseTime: time.Now()})
			if dbc := database.DB.Save(&shortlink); dbc.Error != nil {
				fmt.Println(dbc.Error)
			}

			c.Redirect(http.StatusMovedPermanently, shortlink.Full)
		}
	}
}
