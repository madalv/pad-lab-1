package api

import (
	"context"
	mock_api "course/mocks"
	"course/model"
	"course/pb"
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func getRandomPort() string {
	port := rand.Intn(10000) + 40000
	return fmt.Sprintf(":%d", port)
}

func TestCreateCourse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCourseSvc := mock_api.NewMockcourseService(ctrl)
	mockChapterSvc := mock_api.NewMockchapterService(ctrl)

	listener, err := net.Listen("tcp", ":50100")
	require.NoError(t, err)

	server, err := NewGRPCServer(mockCourseSvc, mockChapterSvc)
	require.NoError(t, err)

	go func() {
		err := server.Serve(listener)
		require.NoError(t, err)
	}()

	conn, err := grpc.Dial(":50100", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewCourseServiceClient(conn)

	createCourse := model.Course{
		Title:       "Very Nice Title",
		AuthorID:    "author-id",
		Description: "This is a Description",
	}

	categoryIds := []string{"1", "2"}

	expectedResponseId := "ididid"

	mockCourseSvc.EXPECT().Create(createCourse, categoryIds).Return(expectedResponseId, nil)

	pbCreateCourseRequest := &pb.CreateCourseRequest{
		Title:       createCourse.Title,
		AuthorId:    createCourse.AuthorID,
		Description: createCourse.Description,
		CategoryIds: categoryIds,
	}

	pbResponseID, err := client.CreateCourse(context.Background(), pbCreateCourseRequest)
	require.NoError(t, err)

	assert.Equal(t, pbResponseID.Id, expectedResponseId)
}

func TestGetCourse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCourseSvc := mock_api.NewMockcourseService(ctrl)
	mockChapterSvc := mock_api.NewMockchapterService(ctrl)

	port := getRandomPort()

	listener, err := net.Listen("tcp", port)
	require.NoError(t, err)

	server, err := NewGRPCServer(mockCourseSvc, mockChapterSvc)
	require.NoError(t, err)

	go func() {
		err := server.Serve(listener)
		require.NoError(t, err)
	}()

	conn, err := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewCourseServiceClient(conn)

	courseID := "course-id"

	expectedCourse := model.Course{
		Base: model.Base{
			ID:        "course-id",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		AuthorID:    "author-id",
		Title:       "title",
		Description: "description",
		Chapters:    []model.Chapter{},
		Categories:  []*model.Category{},
	}

	mockCourseSvc.EXPECT().GetByID(courseID).Return(expectedCourse, nil)

	pbGetCourseRequest := &pb.CourseId{
		Id: courseID,
	}

	pbResponseCourse, err := client.GetCourse(context.Background(), pbGetCourseRequest)
	require.NoError(t, err)

	chapters := make([]model.Chapter, len(pbResponseCourse.Chapters))
	for i := range chapters {
		ch := pbResponseCourse.Chapters[i]
		chapters[i] = model.Chapter{
			Base: model.Base{
				ID: ch.Id,
			},
			Title: ch.Title,
		}
	}

	assert.Equal(t, pbResponseCourse.Course.Id, expectedCourse.ID)
	assert.Equal(t, pbResponseCourse.Course.Title, expectedCourse.Title)
	assert.Equal(t, chapters, expectedCourse.Chapters)
}

func TestGetCourse_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCourseSvc := mock_api.NewMockcourseService(ctrl)
	mockChapterSvc := mock_api.NewMockchapterService(ctrl)

	port := getRandomPort()

	listener, err := net.Listen("tcp", port)
	require.NoError(t, err)

	server, err := NewGRPCServer(mockCourseSvc, mockChapterSvc)
	require.NoError(t, err)

	go func() {
		err := server.Serve(listener)
		require.NoError(t, err)
	}()

	conn, err := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewCourseServiceClient(conn)

	courseID := "course-id"

	mockCourseSvc.EXPECT().GetByID(courseID).Return(model.Course{}, gorm.ErrRecordNotFound)

	pbGetCourseRequest := &pb.CourseId{
		Id: courseID,
	}

	_, err = client.GetCourse(context.Background(), pbGetCourseRequest)

	expectedError := status.Error(codes.Unknown, "record not found")
	assert.EqualError(t, err, expectedError.Error())
}

func TestCreateChapter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCourseSvc := mock_api.NewMockcourseService(ctrl)
	mockChapterSvc := mock_api.NewMockchapterService(ctrl)

	port := getRandomPort()

	listener, err := net.Listen("tcp", port)
	require.NoError(t, err)

	server, err := NewGRPCServer(mockCourseSvc, mockChapterSvc)
	require.NoError(t, err)

	go func() {
		err := server.Serve(listener)
		require.NoError(t, err)
	}()

	conn, err := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewCourseServiceClient(conn)

	createChapter := model.Chapter{
		Title:    "Very Nice Title",
		CourseID: "course-id",
		Body:     "This is a body",
	}

	expectedResponseId := "ididid"

	mockChapterSvc.EXPECT().Create(createChapter).Return(expectedResponseId, nil)

	pbCreateChapterRequest := &pb.CreateChapterRequest{
		Title:    createChapter.Title,
		Body:     createChapter.Body,
		CourseId: createChapter.CourseID,
	}

	pbResponseID, err := client.CreateChapter(context.Background(), pbCreateChapterRequest)
	require.NoError(t, err)

	assert.Equal(t, pbResponseID.Id, expectedResponseId)
}

func TestGetChapter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCourseSvc := mock_api.NewMockcourseService(ctrl)
	mockChapterSvc := mock_api.NewMockchapterService(ctrl)

	port := getRandomPort()
	
	listener, err := net.Listen("tcp", port)
	require.NoError(t, err)

	server, err := NewGRPCServer(mockCourseSvc, mockChapterSvc)
	require.NoError(t, err)

	go func() {
		err := server.Serve(listener)
		require.NoError(t, err)
	}()

	conn, err := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewCourseServiceClient(conn)

	chapterID := "chapter-id"

	expctedChapter := model.Chapter{
		Base: model.Base{
			ID:        chapterID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		CourseID: "course-id",
		Title:    "title",
		Body:     "body body",
	}

	mockChapterSvc.EXPECT().GetByID(chapterID).Return(expctedChapter, nil)

	pbGetChapterRequest := &pb.ChapterId{
		Id: chapterID,
	}

	pbResponseChapter, err := client.GetChapter(context.Background(), pbGetChapterRequest)
	require.NoError(t, err)

	assert.Equal(t, pbResponseChapter.Id, expctedChapter.ID)
	assert.Equal(t, pbResponseChapter.Title, expctedChapter.Title)
	assert.Equal(t, pbResponseChapter.Body, expctedChapter.Body)
}

func TestGetCourseIdsForUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCourseSvc := mock_api.NewMockcourseService(ctrl)
	mockChapterSvc := mock_api.NewMockchapterService(ctrl)

	port := getRandomPort()

	listener, err := net.Listen("tcp", port)
	require.NoError(t, err)

	server, err := NewGRPCServer(mockCourseSvc, mockChapterSvc)
	require.NoError(t, err)

	go func() {
		err := server.Serve(listener)
		require.NoError(t, err)
	}()

	conn, err := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewCourseServiceClient(conn)

	userID := "user-id"

	expectedIDs := []string{"id1", "id2", "id3"}

	mockCourseSvc.EXPECT().GetCourseIDsForUser(userID).Return(expectedIDs)

	pbGetCourseIdsRequest := &pb.UserId{
		Id: userID,
	}

	pbResponseCourseIds, err := client.GetCourseIdsForUser(context.Background(), pbGetCourseIdsRequest)
	require.NoError(t, err)

	ids := make([]string, len(pbResponseCourseIds.Ids))
	for i := range ids {
		id := pbResponseCourseIds.Ids[i]
		ids[i] = id
	}

	assert.Equal(t, expectedIDs, ids)
}
