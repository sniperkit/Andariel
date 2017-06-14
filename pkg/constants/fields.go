/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Inc.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

/*
 * Revision History:
 *     Initial: 10/05/2017        Jia Chenhui
 */

package constants

const (
	// 按条件搜索 GitHub 库时的 query 类型
	QueryLanguage = "language"
	QueryCreated  = "created"

	// 搜索 GitHub 库时的指定语言类型(collection 名称)
	LangC      = "C"
	LangCSharp = "C#"
	LangCPlus  = "C++"
	LangCSS    = "CSS"
	LangGo     = "Go"
	LangHTML   = "HTML"
	LangJava   = "Java"
	LangJS     = "JavaScript"
	LangLua    = "Lua"
	LangObjC   = "Objective-C"
	LangPHP    = "PHP"
	LangPython = "Python"
	LangR      = "R"
	LangRuby   = "Ruby"
	LangScala  = "Scala"
	LangShell  = "Shell"
	LangSwift  = "Swift"

	// 对 GitHub 库搜索结果的排序方式
	SortByStars   = "stars"
	SortByForks   = "forks"
	SortByUpdated = "updated"

	// 对 GitHub 库搜索结果的排序顺序(增/减)
	OrderByAsc  = "asc"
	OrderByDesc = "desc"

	// 搜索库时指定的时间增量
	OneQuarter = "quarter"
	OneMonth   = "month"
	OneWeek    = "week"
	OneDay     = "day"
)