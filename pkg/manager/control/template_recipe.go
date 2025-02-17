package control

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/jaehoonkim/sentinel/pkg/manager/database/vanilla/excute"
	"github.com/jaehoonkim/sentinel/pkg/manager/database/vanilla/stmt"
	"github.com/jaehoonkim/sentinel/pkg/manager/macro/echoutil"
	recipev2 "github.com/jaehoonkim/sentinel/pkg/manager/model/template_recipe/v2"
	"github.com/jaehoonkim/sentinel/pkg/manager/status/state"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// @Description Find TemplateRecipe
// @Security    XAuthToken
// @Produce     json
// @Tags        manager/template_recipe
// @Router      /manager/template_recipe [get]
// @Param       method       query  string false "Template Command Method"
// @Success     200 {array} v2.HttpRsp_TemplateRecipe
func (ctl ControlVanilla) FindTemplateRecipe(ctx echo.Context) (err error) {
	method := echoutil.QueryParam(ctx)["method"]
	buff := bytes.Buffer{}
	for i, s := range strings.Split(method, ".") {
		if 0 < i {
			buff.WriteString(".")
		}
		buff.WriteString(s)
	}
	//뒤에 like 조회 와일드 카드를 붙여준다
	buff.WriteString("%")

	var p = stmt.Limit(50, 1)
	if 0 < len(echoutil.QueryParam(ctx)["p"]) {
		p, err = stmt.PaginationLexer.Parse(echoutil.QueryParam(ctx)["p"])
		err = errors.Wrapf(err, "failed to parse pagination")
		if err != nil {
			return HttpError(err, http.StatusBadRequest)
		}
	}

	rsp := make([]recipev2.HttpRsp_TemplateRecipe, 0, state.ENV__INIT_SLICE_CAPACITY__())

	recipe := recipev2.TemplateRecipe{}
	recipe.Method = buff.String()
	like_method := stmt.Like("method", recipe.Method)
	order := stmt.Asc("name", "args")

	err = ctl.dialect.QueryRows(recipe.TableName(), recipe.ColumnNames(), like_method, order, p)(
		ctx.Request().Context(), ctl)(
		func(scan excute.Scanner, _ int) error {
			err = recipe.Scan(scan)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}

			rsp = append(rsp, recipe)

			return err
		})
	if err != nil {
		return
	}

	return ctx.JSON(http.StatusOK, []recipev2.HttpRsp_TemplateRecipe(rsp))
}
