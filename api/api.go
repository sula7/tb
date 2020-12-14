package api

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"tb/storage"
	"tb/utils"
)

type iTbAPIv1 interface {
	Start()

	ping(c echo.Context) error
	createTagContent(c echo.Context) error
	getTagContentList(c echo.Context) error
	getTagContent(c echo.Context) error
}

type TbAPIv1 struct {
	store storage.ITbStorage
	echo  *echo.Echo
	iTbAPIv1
}

func New(store storage.ITbStorage) *TbAPIv1 {
	e := echo.New()

	api := TbAPIv1{store: store, echo: e}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", api.ping)

	g := e.Group("/api/v1")
	g.POST("/tag-content", api.createTagContent)
	g.GET("/tag-content", api.getTagContentList)
	g.GET("/tag-content/:id", api.getTagContent)

	return &api
}

func (api *TbAPIv1) Start() {
	api.echo.Logger.Fatal(api.echo.Start(":1323"))
}

func (api *TbAPIv1) ping(c echo.Context) error {
	return c.String(http.StatusOK, "ping OK")
}

func (api *TbAPIv1) createTagContent(c echo.Context) error {
	tagContent := TagContent{}

	err := c.Bind(&tagContent)
	if err != nil {
		log.Error("error reading request body:", err)
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
		})
	}

	newsItems, err := utils.GetTagContent(tagContent.URL, tagContent.TagName)
	if err != nil {
		log.Error("error getting & parse tag:", err)
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
		})
	}

	err = api.store.InsertNewsItem(newsItems)
	if err != nil {
		log.Error("error creating tag data into DB:", err)
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "OK",
	})
}

func (api *TbAPIv1) getTagContentList(c echo.Context) error {
	filter := storage.TagContentFilter{
		ContentKeyword: c.QueryParam("tag_content"),
		URL:            c.QueryParam("url"),
		CreatedAtFrom:  c.QueryParam("created_at_from"),
		CreatedAtTo:    c.QueryParam("created_at_to"),
	}

	tagContentList, err := api.store.GetTagContentList(filter)
	if err != nil {
		log.Error("error getting list of tag content:", err)
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "OK",
		Data:    tagContentList,
	})
}

func (api *TbAPIv1) getTagContent(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Error("error getting id:", err)
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: err.Error(),
		})
	}

	tagContent, err := api.store.GetTagContent(id)
	if err != nil {
		log.Error("error getting tag content:", err)
		return c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "OK",
		Data:    tagContent,
	})
}
