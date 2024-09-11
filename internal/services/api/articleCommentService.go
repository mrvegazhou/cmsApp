package api

import (
	"cmsApp/internal/constant"
	"cmsApp/internal/dao"
	"cmsApp/internal/errorx"
	"cmsApp/internal/models"
	"cmsApp/pkg/ip"
	"cmsApp/pkg/utils/arrayx"
	HTMLX "cmsApp/pkg/utils/htmlx"
	"cmsApp/pkg/utils/timex"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

type apiArticleCommentService struct {
	Dao             *dao.AppArticleCommentDao
	ArticleDao      *dao.AppArticleDao
	ArticleReplyDao *dao.AppArticleReplyDao
}

var (
	instanceApiArticleCommentService *apiArticleCommentService
	onceApiArticleCommentService     sync.Once
)

func NewApiArticleCommentService() *apiArticleCommentService {
	onceApiArticleCommentService.Do(func() {
		instanceApiArticleCommentService = &apiArticleCommentService{
			Dao:             dao.NewAppArticleCommentDao(),
			ArticleDao:      dao.NewAppArticleDao(),
			ArticleReplyDao: dao.NewAppArticleReplyDao(),
		}
	})
	return instanceApiArticleCommentService
}

// 保存文章评论
// todo: 需要验证内容是否合法 bert
func (ser *apiArticleCommentService) SaveArticleComment(userId uint64, req models.ArticleCommentPost, ipStr string) (commentInfo models.ArticleCommentWithUserInfo, err error) {

	if &userId == nil {
		return commentInfo, errors.New(constant.ARTICLE_COMMENT_USER_EMPTY_ERR)
	}
	if len(req.Content) == 0 {
		return commentInfo, errors.New(constant.ARTICLE_COMMENT_CONTENT_EMPTY_ERR)
	}
	if utf8.RuneCountInString(req.Content) > constant.ARTICLE_COMMENT_CONTENT_LEN_LIMIT {
		return commentInfo, errors.New(fmt.Sprintf(constant.ARTICLE_COMMENT_CONTENT_LEN_ERR, constant.ARTICLE_COMMENT_CONTENT_LEN_LIMIT))
	}
	// 检查articleId有效性
	articleInfo, err := ser.ArticleDao.GetAppArticle(map[string]interface{}{
		"id": req.ArticleId,
	})
	if err != nil {
		return commentInfo, errors.New(constant.ARTICLE_COMMENT_ARTICLE_ERR)
	}

	comment := models.AppArticleComment{}
	copier.CopyWithOption(&comment, req, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	comment.CreateTime = time.Now()
	comment.ArticleId = req.ArticleId
	// 如果内容包含图片，修改图片为正式存储，删除临时图片信息
	newContent, err := NewApiImgsService().HandleImgs(comment.Content)
	if err != nil {
		return commentInfo, errors.New(constant.ARTICLE_COMMENT_IMG_SRC_ERR)
	}
	comment.Content = comment.Content
	comment.DislikeCount = 0
	comment.LikeCount = 0
	comment.Status = 1
	comment.UserId = userId
	ipInfo := &ip.IPInfo{}
	ipInfo, err = ip.HandleIPInfo(ipStr, "cn")
	comment.Ip = ipInfo.City
	id, err := ser.Dao.CreateArticleComment(comment)
	if err != nil {
		return commentInfo, errors.New(constant.ARTICLE_SAVE_ERR)
	}
	comment.Id = id
	copier.CopyWithOption(&commentInfo, comment, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	// 获取发表人信息
	userInfo, err := NewApiUserService().GetUserInfoRes(map[string]interface{}{"id": comment.UserId})
	commentInfo.UserInfo = userInfo
	commentInfo.Content = newContent
	// 更新文章评论总数
	go func(articleId uint64) {
		NewApiArticleService().UpdateArticleCommentCount(articleId)
	}(articleInfo.Id)
	return commentInfo, nil
}

// todo: 使用bert过滤req.Content
func (ser *apiArticleCommentService) SaveArticleReply(userId uint64, req models.ArticleReplyPost, ipStr string) (reply models.ArticleReplyWithUserAndToReplyContent, err error) {
	copier.CopyWithOption(&reply, req, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	if &userId == nil {
		return reply, errors.New(constant.ARTICLE_REPLY_USER_EMPTY_ERR)
	}
	if &req.CommentId == nil {
		return reply, errors.New(constant.ARTICLE_REPLY_COMMENT_ID_ERR)
	}
	if len(req.Content) == 0 {
		return reply, errors.New(constant.ARTICLE_REPLY_CONTENT_EMPTY_ERR)
	}
	if utf8.RuneCountInString(req.Content) > constant.ARTICLE_COMMENT_CONTENT_LEN_LIMIT {
		return reply, errors.New(fmt.Sprintf(constant.ARTICLE_REPLY_CONTENT_LEN_ERR, constant.ARTICLE_COMMENT_CONTENT_LEN_LIMIT))
	}
	var commentId uint64
	ipInfo := &ip.IPInfo{}
	ipInfo, err = ip.HandleIPInfo(ipStr, "cn")

	// 如果内容包含图片，修改图片为正式存储，删除临时图片信息
	newContent, err := NewApiImgsService().HandleImgs(req.Content)
	if err != nil {
		return reply, errors.New(constant.ARTICLE_COMMENT_IMG_SRC_ERR)
	}
	var commentInfo models.AppArticleComment
	replyModel := models.AppArticleReply{}
	copier.CopyWithOption(&replyModel, req, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	var fromUid, toUid uint64
	fromUid = userId
	if req.ReplyId == 0 {
		// 检查commentId有效性
		commentInfo, err = ser.Dao.GetArticleComment(map[string]interface{}{
			"id": req.CommentId,
		})
		if err != nil {
			return reply, errors.New(constant.ARTICLE_REPLY_COMMENT_ID_ERR)
		}
		reply.ToReplyContent = commentInfo.Content

		replyModel.ArticleId = commentInfo.ArticleId
		replyModel.CommentId = commentInfo.Id
		replyModel.LikeCount = 0
		replyModel.DislikeCount = 0
		replyModel.CreateTime = time.Now()
		replyModel.FromUid = userId
		replyModel.ToUid = commentInfo.UserId
		replyModel.Pids = cast.ToString(commentInfo.Id)
		replyModel.Ip = ipInfo.City
		replyModel.Content = req.Content
		commentId = commentInfo.Id
		toUid = commentInfo.UserId
	} else if req.ReplyId != 0 {
		// 是评论的评论
		replyInfo, err := ser.ArticleReplyDao.GetArticleReply(map[string]interface{}{
			"id": req.ReplyId,
		})
		if err != nil {
			return reply, errors.New(constant.ARTICLE_REPLY_COMMENT_ID_ERR)
		}
		reply.ToReplyContent = replyInfo.Content

		replyModel.ArticleId = replyInfo.ArticleId
		replyModel.CommentId = replyInfo.CommentId
		replyModel.ReplyId = replyInfo.Id
		replyModel.LikeCount = 0
		replyModel.DislikeCount = 0
		replyModel.CreateTime = time.Now()
		replyModel.FromUid = userId
		replyModel.ToUid = replyInfo.FromUid
		replyModel.Pids = fmt.Sprintf("%s,%s", replyInfo.Pids, replyInfo.Id)
		replyModel.Ip = ipInfo.City
		replyModel.Content = req.Content
		commentId = replyInfo.CommentId
		toUid = replyInfo.FromUid
	}
	replyId, err := ser.ArticleReplyDao.CreateArticleReply(replyModel)
	if err != nil {
		return reply, errors.New(constant.ARTICLE_REPLY_SAVE_ERR)
	}

	reply.Id = replyId
	copier.CopyWithOption(&reply, replyModel, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	// 获取评论人 和 被评论人的信息
	uids := [2]uint64{fromUid, toUid}
	userList, err := NewApiUserService().GetUserList(uids[:2])
	if err != nil {
		return reply, errors.New(constant.ARTICLE_COMMENT_REPLY_USERS_ERR)
	}
	usersMap := make(map[uint64]models.AppUserInfo)
	for _, user := range userList {
		userModel := models.AppUserInfo{}
		copier.CopyWithOption(&userModel, user, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		usersMap[user.Id] = userModel
	}
	reply.FromUser = usersMap[fromUid]
	reply.ToUser = usersMap[toUid]
	reply.Content = newContent
	// 更新评论表的reply top 和 文章的评论总数
	go func(commentId uint64, commentInfo models.AppArticleComment) {
		ser.UpdateCommentTopRepies(commentId, replyId, commentInfo)
		NewApiArticleService().UpdateArticleCommentCount(commentInfo.ArticleId)
	}(commentId, commentInfo)

	return reply, nil
}

// 更新top reply
func (ser *apiArticleCommentService) UpdateCommentTopRepies(commentId, replyId uint64, commentInfo models.AppArticleComment) {
	var err error
	if &commentInfo == nil {
		commentInfo, err = ser.Dao.GetArticleComment(map[string]interface{}{
			"id": commentId,
		})
	}
	if err == nil {
		topReplyIds := strings.Split(commentInfo.TopReplyIds, ",")
		if len(topReplyIds) < 2 {
			topReplyIds = append(topReplyIds, cast.ToString(replyId))
		} else {
			// 当长度等于2时，移除第一个元素，追加新元素
			topReplyIds[0] = cast.ToString(replyId)
			topReplyIds = topReplyIds[:2]
		}
		commentInfo.TopReplyIds = strings.Join(topReplyIds, ",")
		commentInfo.ReplyCount = commentInfo.ReplyCount + 1
		ser.Dao.UpdateArticleComment(commentInfo)
	}
}

func (ser *apiArticleCommentService) GetCommentList(req models.ArticleCommentListPost) (commentReplies models.ArticleCommentListResp, err error) {
	if &req.ArticleId == nil {
		return commentReplies, errors.New(constant.ARTICLE_COMMENT_ARTICLE_ERR)
	}
	// 获取文章
	_, err = ser.ArticleDao.GetAppArticle(map[string]interface{}{
		"id": req.ArticleId,
	})
	if err != nil {
		return commentReplies, errors.New(constant.ARTICLE_COMMENT_ARTICLE_ERR)
	}

	var orderBy string
	if req.OrderBy == "" {
		orderBy = "create_time desc"
	} else {
		orderBy = fmt.Sprintf("%s desc, id desc", req.OrderBy)
	}

	if !timex.IsValidTimestamp(req.CurrentTime) {
		return commentReplies, errors.New(constant.ARTICLE_COMMENT_CURRENT_TIME_ERR)
	}
	pageTime := time.Unix(req.CurrentTime, 0)

	conditions := map[string][]interface{}{
		"article_id":  {"= ?", req.ArticleId},
		"create_time": {"<= ?", pageTime},
		"Status":      {"= ?", 1},
	}

	listComment, _, err := ser.Dao.GetArticleCommentListNoTotal(conditions, req.Page, constant.ARTICLE_COMMENT_PAGE_SIZE, orderBy)
	commentMap := make(map[uint64]models.AppArticleComment)

	for _, comment := range listComment {
		newContent := HTMLX.AppendParamToImageSrc(comment.Content)
		comment.Content = newContent
		commentMap[comment.Id] = comment
	}

	commentReplies, err = ser.GetRecentReplies(commentMap, orderBy)
	if err != nil {
		return commentReplies, errorx.NewCustomError(errorx.HTTP_UNKNOW_ERR, constant.ARTICLE_COMMENT_PAGE_ERR)
	}
	// 判断HasNext
	if len(listComment) < constant.ARTICLE_COMMENT_PAGE_SIZE {
		commentReplies.HasNext = false
	} else {
		commentReplies.HasNext = true
	}
	return commentReplies, nil
}

// 获取评论列表中的回复列表 2条
// todo: 这里可以加缓存，缓存key可以是升序排序的评论id
func (ser *apiArticleCommentService) GetRecentReplies(commentMap map[uint64]models.AppArticleComment, orderBy string) (commentReplies models.ArticleCommentListResp, err error) {
	//
	if len(commentMap) == 0 {
		return commentReplies, nil
	}
	var replyIds []uint64
	var userIds []uint64

	for _, value := range commentMap {
		ids := strings.Split(value.TopReplyIds, ",")
		for _, replyId := range ids {
			replyIds = append(replyIds, cast.ToUint64(replyId))
		}
		userIds = append(userIds, value.UserId)
	}

	listReplies, err := ser.GetReplyByIn(replyIds, orderBy)
	if err != nil {
		return commentReplies, err
	}

	var pReplyIds []uint64
	for _, replyInfo := range listReplies {
		userIds = append(userIds, replyInfo.ToUid, replyInfo.ToUid)
		if replyInfo.ReplyId != 0 {
			pReplyIds = append(pReplyIds, replyInfo.ReplyId)
		}
	}
	pReplies, err := ser.GetReplyByIn(pReplyIds, orderBy)
	preplyMap := make(map[uint64]models.AppArticleReply)
	for _, pReply := range pReplies {
		preplyMap[pReply.Id] = pReply
	}

	userMap, err := NewApiUserService().GetUserMapListByIds(userIds)
	if err != nil {
		return commentReplies, err
	}

	commentWithReplyList := make(map[uint64][]models.ArticleReplyWithUserAndToReplyContent)
	// 循环reply list 添加user info
	for _, replyInfo := range listReplies {
		replyModel := models.ArticleReplyWithUserAndToReplyContent{}
		copier.CopyWithOption(&replyModel, replyInfo, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		replyModel.ToUser = userMap[replyInfo.ToUid]
		replyModel.FromUser = userMap[replyInfo.FromUid]
		replyModel.ToReplyContent = ""
		if replyInfo.ReplyId != 0 {
			replyModel.ToReplyContent = preplyMap[replyInfo.ReplyId].Content
		}
		commentWithReplyList[replyInfo.CommentId] = append(commentWithReplyList[replyInfo.CommentId], replyModel)
	}

	// 组装成ArticleCommentListResp
	for id, comment := range commentMap {
		model := models.ArticleCommentReplies{}
		model.Comment.UserInfo = userMap[comment.UserId]
		commentModel := models.ArticleCommentResp{}
		copier.CopyWithOption(&commentModel, comment, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		model.Comment.ArticleCommentResp = commentModel
		model.Replies = commentWithReplyList[id]
		model.HasNext = false
		if comment.ReplyCount > 2 {
			model.HasNext = true
		}
		commentReplies.CommentList = append(commentReplies.CommentList, model)
	}
	sort.Slice(commentReplies.CommentList, func(i, j int) bool {
		return commentReplies.CommentList[j].Comment.CreateTime.Before(commentReplies.CommentList[i].Comment.CreateTime) // 升序
	})
	commentReplies.HasNext = false
	return commentReplies, nil
}

// 获取评论的回复列表
func (ser *apiArticleCommentService) GetReplyList(req models.ArticleReplyListPost) (listReplyWithContent models.ArticleReplyListResp, err error) {
	if &req.CommentId == nil {
		return listReplyWithContent, errorx.NewCustomError(errorx.HTTP_BIND_PARAMS_ERR, constant.ARTICLE_REPLY_COMMENT_ID_ERR)
	}
	seconds := req.CurrentTime / 1000
	if !timex.IsValidTimestamp(seconds) {
		return listReplyWithContent, errorx.NewCustomError(errorx.HTTP_BIND_PARAMS_ERR, constant.ARTICLE_COMMENT_PAGE_ERR)
	}
	pageTime := time.Unix(seconds, 0)

	var orderBy string
	if req.OrderBy == "" {
		orderBy = "create_time desc"
	} else {
		orderBy = fmt.Sprintf("%s desc, id desc", req.OrderBy)
	}

	// 获取父级帖子的评论总数
	var count uint
	var conditions map[string][]interface{}
	if req.ReplyId != 0 {
		replyInfo, err := ser.ArticleReplyDao.GetArticleReply(map[string]interface{}{
			"id": req.ReplyId,
		})
		if err != nil {
			return listReplyWithContent, errorx.NewCustomError(errorx.HTTP_BIND_PARAMS_ERR, constant.ARTICLE_REPLY_ID_ERR)
		}
		count = replyInfo.ReplyCount
		conditions = map[string][]interface{}{
			"reply_id":    {"= ?", req.ReplyId},
			"create_time": {"<= ?", pageTime},
			"Status":      {"= ?", 1},
		}
	} else {
		commentInfo, err := ser.Dao.GetArticleComment(map[string]interface{}{
			"id": req.CommentId,
		})
		if err != nil {
			return listReplyWithContent, errorx.NewCustomError(errorx.HTTP_BIND_PARAMS_ERR, constant.ARTICLE_REPLY_COMMENT_ID_ERR)
		}
		count = commentInfo.ReplyCount
		conditions = map[string][]interface{}{
			"comment_id":  {"= ?", req.CommentId},
			"create_time": {"<= ?", pageTime},
			"Status":      {"= ?", 1},
		}
	}
	listReplyWithContent.HasNext = false
	offset := (req.Page - 1) * constant.ARTICLE_COMMENT_PAGE_SIZE
	if count > cast.ToUint(offset+constant.ARTICLE_COMMENT_PAGE_SIZE) {
		listReplyWithContent.HasNext = true
	}

	listReply, err := ser.ArticleReplyDao.GetArticleReplyListNoTotal(conditions, req.Page, constant.ARTICLE_COMMENT_PAGE_SIZE, orderBy)

	var userIds []uint64
	var pReplyIds []uint64
	for _, reply := range listReply {
		userIds = append(userIds, reply.ToUid, reply.FromUid)
		if reply.ReplyId != 0 {
			pReplyIds = append(pReplyIds, reply.ReplyId)
		}
	}
	userMap, err := NewApiUserService().GetUserMapListByIds(userIds)
	// 回复的帖子
	pReplies, err := ser.GetReplyByIn(pReplyIds, orderBy)
	preplyMap := make(map[uint64]models.AppArticleReply)
	for _, pReply := range pReplies {
		preplyMap[pReply.Id] = pReply
	}

	list := []models.ArticleReplyWithUserAndToReplyContent{}
	for _, reply := range listReply {
		model := models.ArticleReplyWithUserAndToReplyContent{}
		copier.CopyWithOption(&model, reply, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		if reply.ReplyId != 0 {
			model.ToReplyContent = preplyMap[reply.ReplyId].Content
		} else {
			model.ToReplyContent = ""
		}
		model.ToUser = userMap[reply.ToUid]
		model.FromUser = userMap[reply.FromUid]
		list = append(list, model)
	}
	listReplyWithContent.ReplyList = list
	listReplyWithContent.CurrentTime = req.CurrentTime
	listReplyWithContent.Page = req.Page
	return
}

// 通过in条件查询回复列表
func (ser *apiArticleCommentService) GetReplyByIn(replyIds []uint64, orderBy string) (replies []models.AppArticleReply, err error) {
	if len(replyIds) == 0 {
		return
	}
	replyIds = arrayx.RemoveRepeatedElement(replyIds) //回复id去重
	conditions := map[string][]interface{}{
		"id": {"in (?)", replyIds},
	}
	replies, err = ser.ArticleReplyDao.GetArticleReplies(conditions, orderBy)
	if err != nil {
		return replies, err
	}
	return
}

func (ser *apiArticleCommentService) HandleReport(req models.ArticleCommentReportReq) (success bool, err error) {
	return true, nil
}
