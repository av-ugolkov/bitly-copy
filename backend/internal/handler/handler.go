package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/av-ugolkov/bitly-copy/pkg/url"
	"github.com/labstack/echo/v4"
)

type service interface {
	Get(ctx context.Context, url string) (string, error)
	Add(ctx context.Context, uid, url, aliase string) (string, error)
	GetUrls(ctx context.Context, uid string) ([]string, error)
	Remove(ctx context.Context, uid, url string) error
	Statistics(ctx context.Context, shortUrl string) (string, error)
}

type handler struct {
	svc service
}

func Create(r *echo.Echo, svc service) {
	h := handler{
		svc: svc,
	}

	r.GET("/:shortURL", h.redirect)
	r.POST("/add", h.add)
	r.GET("/get-urls", h.getUrls)
	r.DELETE("/remove", h.remove)
	r.GET("/statistics", h.statistics)
}

func (h *handler) redirect(c echo.Context) error {
	ctx := c.Request().Context()
	url, err := h.svc.Get(ctx, c.Request().URL.Path[1:])
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.Redirect(http.StatusFound, url)
}

func (h *handler) add(c echo.Context) error {
	var request struct {
		UID    string `json:"uid"`
		URL    string `json:"url"`
		Aliase string `json:"aliase"`
	}

	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	if !url.Validate(request.URL) {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "invalid url",
		})
	}

	if request.Aliase == "statistics" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "invalid aliase",
		})
	}

	ctx := c.Request().Context()
	shortURL, err := h.svc.Add(ctx, request.UID, request.URL, request.Aliase)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"short_url": fmt.Sprintf("http://%s/%s", c.Request().Host, shortURL),
	})
}

func (h *handler) getUrls(c echo.Context) error {
	ctx := c.Request().Context()

	uid := c.QueryParam("uid")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "uid is required",
		})
	}

	urls, err := h.svc.GetUrls(ctx, uid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	shortUrls := make([]string, 0, len(urls))
	for _, url := range urls {
		shortUrls = append(shortUrls, fmt.Sprintf("http://%s/%s", c.Request().Host, url))
	}

	return c.JSON(http.StatusOK, echo.Map{
		"urls": shortUrls,
	})
}

func (h *handler) remove(c echo.Context) error {
	ctx := c.Request().Context()

	var request struct {
		UID      string `json:"uid"`
		ShortURL string `json:"short_url"`
	}

	err := c.Bind(&request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	err = h.svc.Remove(ctx, request.UID, request.ShortURL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, nil)
}

func (h *handler) statistics(c echo.Context) error {
	ctx := c.Request().Context()

	shortUrl := c.QueryParam("short_url")
	if shortUrl == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "short_url is required",
		})
	}

	statistic, err := h.svc.Statistics(ctx, shortUrl)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"short_url":  shortUrl,
		"statistics": statistic,
	})
}
