// Copyright 2017 Vckai Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package admin

import (
	"math"

	"github.com/vckai/novel/app/models"
	"github.com/vckai/novel/app/services"
	"github.com/vckai/novel/app/utils"
	"github.com/vckai/novel/app/utils/log"
)

type ChapterController struct {
	BaseController
}

// 章节列表
func (this *ChapterController) Index() {
	q := this.GetString("q")
	novId, _ := this.GetUint32("novid")
	p, _ := this.GetInt("p", 1)
	if novId < 1 {
		this.Msg("参数错误，无法访问")
	}

	size := 10
	offset := (p - 1) * size
	chapters, count := services.ChapterService.GetNovChaps(novId, size, offset, "desc", true)

	search := map[string]interface{}{
		"p":     p,
		"q":     q,
		"novId": novId,
	}
	this.Data["Search"] = search
	this.Data["Chapters"] = chapters
	this.Data["ChaptersCount"] = count
	this.Data["MaxPages"] = math.Ceil(float64(count) / float64(size))
	this.View("chapter/index.tpl")
}

// 添加小说章节页面
func (this *ChapterController) Add() {
	// post 数据提交
	if this.IsAjax() {
		this.save()
		return
	}

	this.Data["Chapter"] = models.NewChapter()
	this.Data["IsUpdate"] = false
	this.Data["PostUrl"] = utils.URLFor("admin.ChapterController.Add")
	this.View("chapter/add.tpl")
}

// 编辑小说页面
func (this *ChapterController) Edit() {
	// post 数据提交
	if this.IsAjax() {
		this.save()
		return
	}

	id, _ := this.GetUint64("id")
	novId, _ := this.GetUint32("novid")
	if id < 1 || novId < 1 {
		this.Msg("参数错误，无法访问")
	}

	chapter := services.ChapterService.Get(id, novId)
	if chapter == nil {
		this.Msg("该小说不存在或者已被删除")
	}

	this.Data["Chapter"] = chapter
	this.Data["IsUpdate"] = true
	this.Data["PostUrl"] = utils.URLFor("admin.ChapterController.Edit")
	this.View("chapter/add.tpl")
}

// 删除小说
func (this *ChapterController) Delete() {
	id, _ := this.GetUint64("id")
	novId, _ := this.GetUint32("novid")
	name := this.GetString("name")
	if id < 1 {
		this.Msg("参数错误，无法访问")
	}

	err := services.ChapterService.Delete(id, novId)
	if err != nil {
		this.OutJson(1001, "删除小说失败")
	}

	// 添加操作日记
	this.AddLog(3003, "删除", name, id)

	this.OutJson(0, "已删除！")
}

// 保存数据
// 提供修改/新增处理
func (this *ChapterController) save() {
	id, _ := this.GetUint64("chapter_id")
	novId, _ := this.GetUint32("novid")

	mtitle := "添加"
	chapter := models.NewChapter()
	// 编辑
	if id > 0 {
		chapter = services.ChapterService.Get(id, novId)
		if chapter == nil {
			this.OutJson(1001, "该章节不存在或者已被删除")
		}

		mtitle = "修改"
	}

	// 入库参数
	chapter.Title = this.GetString("title")
	chapter.Desc = this.GetString("desc")
	chapter.NovId = novId

	err := services.ChapterService.Save(chapter)
	if err != nil {
		log.Error(mtitle, "章节失败：", err.Error())
		this.OutJson(1002, mtitle+"章节失败："+err.Error())
	}

	// 添加操作日记
	this.AddLog(3003, mtitle, chapter.Title, chapter.Id)

	this.OutJson(0, mtitle+"成功")
}
