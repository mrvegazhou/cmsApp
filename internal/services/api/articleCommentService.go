package api

import (
	"cmsApp/internal/constant"
	"cmsApp/internal/dao"
	"cmsApp/internal/models"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/spf13/cast"
	"sync"
	"time"
	"unicode/utf8"
)

type apiArticleCommentService struct {
	Dao               *dao.AppArticleCommentDao
	ArticleDao        *dao.AppArticleDao
	ArticleCommentDao *dao.AppArticleCommentDao
	ArticleReplyDao   *dao.AppArticleReplyDao
}

var (
	instanceApiArticleCommentService *apiArticleCommentService
	onceApiArticleCommentService     sync.Once
)

func NewApiArticleCommentService() *apiArticleCommentService {
	onceApiArticleCommentService.Do(func() {
		instanceApiArticleCommentService = &apiArticleCommentService{
			Dao:               dao.NewAppArticleCommentDao(),
			ArticleDao:        dao.NewAppArticleDao(),
			ArticleCommentDao: dao.NewAppArticleCommentDao(),
			ArticleReplyDao:   dao.NewAppArticleReplyDao(),
		}
	})
	return instanceApiArticleCommentService
}

// 保存文章评论
func (ser *apiArticleCommentService) SaveArticleComment(userId uint64, req models.ArticleCommentPost) (comment models.AppArticleComment, err error) {
	copier.CopyWithOption(&comment, req, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	if &userId == nil {
		return comment, errors.New(constant.ARTICLE_COMMENT_USER_EMPTY_ERR)
	}
	if len(req.Content) == 0 {
		return comment, errors.New(constant.ARTICLE_COMMENT_CONTENT_EMPTY_ERR)
	}
	if utf8.RuneCountInString(req.Content) > constant.ARTICLE_COMMENT_CONTENT_LEN_LIMIT {
		return comment, errors.New(fmt.Sprintf(constant.ARTICLE_COMMENT_CONTENT_LEN_ERR, constant.ARTICLE_COMMENT_CONTENT_LEN_LIMIT))
	}
	// 检查articleId有效性
	_, err = ser.ArticleDao.GetAppArticle(map[string]interface{}{
		"id": req.ArticleId,
	})
	if err != nil {
		return comment, errors.New(constant.ARTICLE_COMMENT_ARTICLE_ERR)
	}
	comment.CreateTime = time.Now()
	comment.ArticleId = req.ArticleId
	comment.Content = req.Content
	comment.DislikeCount = 0
	comment.LikeCount = 0
	comment.Status = 1
	comment.UserId = userId
	id, err := ser.Dao.CreateArticleComment(comment)
	if err != nil {
		return comment, errors.New(constant.ARTICLE_SAVE_ERR)
	}
	comment.Id = id
	return comment, nil
}

func (ser *apiArticleCommentService) SaveArticleReply(userId uint64, req models.ArticleReplyPost) (reply models.AppArticleReply, err error) {
	copier.CopyWithOption(&reply, req, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	if &userId == nil {
		return reply, errors.New(constant.ARTICLE_REPLY_USER_EMPTY_ERR)
	}
	if &req.PostId == nil {
		return reply, errors.New(constant.ARTICLE_REPLY_COMMENT_ID_ERR)
	}
	if len(req.Content) == 0 {
		return reply, errors.New(constant.ARTICLE_REPLY_CONTENT_EMPTY_ERR)
	}
	if utf8.RuneCountInString(req.Content) > constant.ARTICLE_COMMENT_CONTENT_LEN_LIMIT {
		return reply, errors.New(fmt.Sprintf(constant.ARTICLE_REPLY_CONTENT_LEN_ERR, constant.ARTICLE_COMMENT_CONTENT_LEN_LIMIT))
	}
	if req.Type == 1 {
		// 检查commentId有效性
		commentInfo, err := ser.ArticleCommentDao.GetArticleComment(map[string]interface{}{
			"id": req.PostId,
		})
		if err != nil {
			return reply, errors.New(constant.ARTICLE_REPLY_COMMENT_ID_ERR)
		}
		reply.ArticleId = commentInfo.ArticleId
		reply.PostId = commentInfo.Id
		reply.LikeCount = 0
		reply.DislikeCount = 0
		reply.CreateTime = time.Now()
		reply.FromUid = userId
		reply.ToUid = commentInfo.UserId
		reply.Type = 1
		reply.Pids = cast.ToString(commentInfo.Id)
	} else if req.Type == 2 {
		// 是评论的评论
		replyInfo, err := ser.ArticleReplyDao.GetArticleReply(map[string]interface{}{
			"id": req.PostId,
		})
		if err != nil {
			return reply, errors.New(constant.ARTICLE_REPLY_COMMENT_ID_ERR)
		}
		reply.ArticleId = replyInfo.ArticleId
		reply.PostId = replyInfo.Id
		reply.LikeCount = 0
		reply.DislikeCount = 0
		reply.CreateTime = time.Now()
		reply.FromUid = userId
		reply.ToUid = replyInfo.FromUid
		reply.Pids = fmt.Sprintf("%s,%s", replyInfo.Pids, replyInfo.Id)
		reply.Type = 2
	}
	return reply, nil
}

func (ser *apiArticleCommentService) GetCommentList(req models.ArticleCommentListPost) (listComment []models.AppArticleComment, err error) {
	if &req.ArticleId == nil {
		return listComment, errors.New(constant.ARTICLE_COMMENT_ARTICLE_ERR)
	}
	conditions := map[string][]interface{}{
		"article_id": {"= ?", req.ArticleId},
	}
	listComment, _, _, err = ser.ArticleCommentDao.GetArticleCommentList(conditions, req.OffSetPage, constant.ARTICLE_COMMENT_PAGE_SIZE)
	return
}

func (ser *apiArticleCommentService) GetReplyList(req models.ArticleReplyListPost) (listReply []models.AppArticleReply, err error) {
	if &req.CommentId == nil {
		return listReply, errors.New(constant.ARTICLE_REPLY_COMMENT_ID_ERR)
	}
	conditions := map[string][]interface{}{
		"post_id": {"= ?", req.CommentId},
		"type":    {"= ?", 1},
	}

}
