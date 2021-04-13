package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/exam105-UPD/backend/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type questionUsecase struct {
	questionRepo   domain.QuestionRepository
	contextTimeout time.Duration
}

// NewQuestionUsecase will create new an questionUsecase object representation of domain.QuestionUsecase interface in McqCompleteQuestion.go file
func NewQuestionUsecase(qsRepo domain.QuestionRepository, timeout time.Duration) domain.QuestionUsecase {
	return &questionUsecase{
		questionRepo:   qsRepo,
		contextTimeout: timeout,
	}
}

func (qsUC *questionUsecase) SaveMCQ(requestCtx context.Context, allMcqs *domain.MCQModel, username string, useremail string) (error) {
	
	ctx, cancel := context.WithTimeout(requestCtx, qsUC.contextTimeout)
	defer cancel()

	metadataBson := new(domain.MetadataBson)
	questionSet := []interface{}{}
	questionHexIds := make([]string, 0)

	for key, currentQuesPointer := range *allMcqs {

		singleQuestion := new(domain.Question)

		fmt.Printf("\n *************************** \n")
		if key == 0 {

			metadataBson.System = currentQuesPointer.System
			metadataBson.Board = currentQuesPointer.Board
			metadataBson.Subject = currentQuesPointer.Subject
			metadataBson.Year = currentQuesPointer.Year
			metadataBson.Month = currentQuesPointer.Month
			metadataBson.Series = currentQuesPointer.Series
			metadataBson.Paper = currentQuesPointer.Paper
			metadataBson.Username = username
			metadataBson.Useremail = useremail

		} else {
			_id := primitive.NewObjectID()
			questionText := currentQuesPointer.Questions
			marks := currentQuesPointer.Marks
			optionsArray := make([]domain.QuestionOptions, 0)
			topicsArray := make([]domain.QuestionTopics, 0)
			imagesArray := make([]domain.QuestionImages, 0)

			// fmt.Printf("Question:-->  %s --- Marks:--> %s \n ", questionText, marks)
			for _, option := range currentQuesPointer.Options {

				//fmt.Printf("Key---> %d \n", key)
				qsOption := new(domain.QuestionOptions)
				qsOption.Option = option.Option
				qsOption.Correct = option.Correct
				optionsArray = append(optionsArray, *qsOption)

				// fmt.Printf("Option: %s --- Correct: %t \n", qsOption.Option, qsOption.Correct)
				// fmt.Printf("%v \n --------------------", optionsArray)

			}

			for _, topic := range currentQuesPointer.Topics {

				//fmt.Printf("Key---> %d \n", key)
				qsTopic := new(domain.QuestionTopics)
				qsTopic.Topic = topic.Topic
				topicsArray = append(topicsArray, *qsTopic)

				// fmt.Printf("Topic: %s  \n", qsTopic.Topic)
				// fmt.Printf("%v \n --------------------", topicsArray)

			}

			for _, imageurl := range currentQuesPointer.Images {

				//fmt.Printf("Key---> %d \n", key)
				qsImageUrl := new(domain.QuestionImages)
				qsImageUrl.Imageurl = imageurl.Imageurl
				imagesArray = append(imagesArray, *qsImageUrl)

				// fmt.Printf("ImageURL: %s  \n", qsImageUrl.Imageurl)
				// fmt.Printf("%v \n --------------------", imagesArray)

			}

			singleQuestion.ID = _id
			singleQuestion.Questions = questionText
			singleQuestion.Marks = marks
			singleQuestion.Options = optionsArray
			singleQuestion.Topics = topicsArray
			singleQuestion.Images = imagesArray

			questionSet = append(questionSet, singleQuestion)

			//adding hexID to metadata
			questionHexIds = append(questionHexIds, _id.Hex())

			fmt.Printf("Single Qs ->>> %v \t \n", singleQuestion)
		}
	}

	metadataBson.QuestionHexIds = questionHexIds

	fmt.Println("*******************************________________**********************************")
	fmt.Printf("\n %#v", questionSet)
	fmt.Printf("\n ID---- %#v", questionHexIds)

	qsUC.questionRepo.SaveQuestionMetadata(ctx, metadataBson)
	qsUC.questionRepo.SaveAllQuestions(ctx, questionSet)
	return nil
}

func (qsUC *questionUsecase) SaveTheoryQuestion(requestCtx context.Context, allMcqs *domain.TheoryModel, username string, useremail string) (error) {

	ctx, cancel := context.WithTimeout(requestCtx, qsUC.contextTimeout)
	defer cancel()

	metadataBson := new(domain.MetadataBson)
	questionSet := []interface{}{}
	questionHexIds := make([]string, 0)

	for key, currentQuesPointer := range *allMcqs {

		singleQuestion := new(domain.TheoryQuestion)

		fmt.Printf("\n *************************** \n")
		if key == 0 {

			metadataBson.System = currentQuesPointer.System
			metadataBson.Board = currentQuesPointer.Board
			metadataBson.Subject = currentQuesPointer.Subject
			metadataBson.Year = currentQuesPointer.Year
			metadataBson.Month = currentQuesPointer.Month
			metadataBson.Series = currentQuesPointer.Series
			metadataBson.Paper = currentQuesPointer.Paper
			metadataBson.Username = username
			metadataBson.Useremail = useremail

		} else {
			_id := primitive.NewObjectID()
			questionText := currentQuesPointer.Questions
			marks := currentQuesPointer.Marks
			answer := currentQuesPointer.Answer
			topicsArray := make([]domain.QuestionTopics, 0)
			imagesArray := make([]domain.QuestionImages, 0)

			for _, topic := range currentQuesPointer.Topics {

				//fmt.Printf("Key---> %d \n", key)
				qsTopic := new(domain.QuestionTopics)
				qsTopic.Topic = topic.Topic
				topicsArray = append(topicsArray, *qsTopic)

				// fmt.Printf("Topic: %s  \n", qsTopic.Topic)
				// fmt.Printf("%v \n --------------------", topicsArray)

			}

			for _, imageurl := range currentQuesPointer.Images {

				//fmt.Printf("Key---> %d \n", key)
				qsImageUrl := new(domain.QuestionImages)
				qsImageUrl.Imageurl = imageurl.Imageurl
				imagesArray = append(imagesArray, *qsImageUrl)

				// fmt.Printf("ImageURL: %s  \n", qsImageUrl.Imageurl)
				// fmt.Printf("%v \n --------------------", imagesArray)

			}

			singleQuestion.ID = _id
			singleQuestion.Question = questionText
			singleQuestion.Marks = marks
			singleQuestion.Answer = answer
			singleQuestion.Topics = topicsArray
			singleQuestion.Images = imagesArray

			questionSet = append(questionSet, singleQuestion)

			//adding hexID to metadata
			questionHexIds = append(questionHexIds, _id.Hex())

			fmt.Printf("Single Qs ->>> %v \t \n", singleQuestion)
		}
	}

	metadataBson.QuestionHexIds = questionHexIds

	fmt.Println("*******************************________________**********************************")
	fmt.Printf("\n %#v", questionSet)
	fmt.Printf("\n ID---- %#v", questionHexIds)

	qsUC.questionRepo.SaveQuestionMetadata(ctx, metadataBson)
	qsUC.questionRepo.SaveAllQuestions(ctx, questionSet)
	return nil

}

func (qsUC *questionUsecase) AddSingleQuestion(requestCtx context.Context, singleMCQ *domain.Question, metadataID string) (int64, error) {
	
	ctx, cancel := context.WithTimeout(requestCtx, qsUC.contextTimeout)
	defer cancel()

	singleQuestion := new(domain.Question)
	
	_id := primitive.NewObjectID()
	questionText := singleMCQ.Questions
	marks := singleMCQ.Marks
	optionsArray := make([]domain.QuestionOptions, 0)
	topicsArray := make([]domain.QuestionTopics, 0)
	imagesArray := make([]domain.QuestionImages, 0)

	// fmt.Printf("Question:-->  %s --- Marks:--> %s \n ", questionText, marks)
	for _, option := range singleMCQ.Options {

		//fmt.Printf("Key---> %d \n", key)
		qsOption := new(domain.QuestionOptions)
		qsOption.Option = option.Option
		qsOption.Correct = option.Correct
		optionsArray = append(optionsArray, *qsOption)

		// fmt.Printf("Option: %s --- Correct: %t \n", qsOption.Option, qsOption.Correct)
		// fmt.Printf("%v \n --------------------", optionsArray)

	}

	for _, topic := range singleMCQ.Topics {

		//fmt.Printf("Key---> %d \n", key)
		qsTopic := new(domain.QuestionTopics)
		qsTopic.Topic = topic.Topic
		topicsArray = append(topicsArray, *qsTopic)

		// fmt.Printf("Topic: %s  \n", qsTopic.Topic)
		// fmt.Printf("%v \n --------------------", topicsArray)

	}

	for _, imageurl := range singleMCQ.Images {

		//fmt.Printf("Key---> %d \n", key)
		qsImageUrl := new(domain.QuestionImages)
		qsImageUrl.Imageurl = imageurl.Imageurl
		imagesArray = append(imagesArray, *qsImageUrl)

		// fmt.Printf("ImageURL: %s  \n", qsImageUrl.Imageurl)
		// fmt.Printf("%v \n --------------------", imagesArray)

	}

	singleQuestion.ID = _id
	singleQuestion.Questions = questionText
	singleQuestion.Marks = marks
	singleQuestion.Options = optionsArray
	singleQuestion.Topics = topicsArray
	singleQuestion.Images = imagesArray

	fmt.Printf("Single Qs ->>> %v \t \n", singleQuestion)
	return qsUC.questionRepo.AddSingleQuestion(ctx, singleQuestion, metadataID)
	//return &singleQuestion, nil

}

func (qsUC *questionUsecase) AddSingleTheoryQuestion(requestCtx context.Context, singleMCQ *domain.TheoryQuestion, metadataID string) (int64, error) {
	
	ctx, cancel := context.WithTimeout(requestCtx, qsUC.contextTimeout)
	defer cancel()

	singleQuestion := new(domain.TheoryQuestion)
	
	_id := primitive.NewObjectID()
	questionText := singleMCQ.Question
	marks := singleMCQ.Marks
	answer := singleMCQ.Answer
	topicsArray := make([]domain.QuestionTopics, 0)
	imagesArray := make([]domain.QuestionImages, 0)

	for _, topic := range singleMCQ.Topics {

		//fmt.Printf("Key---> %d \n", key)
		qsTopic := new(domain.QuestionTopics)
		qsTopic.Topic = topic.Topic
		topicsArray = append(topicsArray, *qsTopic)

		// fmt.Printf("Topic: %s  \n", qsTopic.Topic)
		// fmt.Printf("%v \n --------------------", topicsArray)

	}

	for _, imageurl := range singleMCQ.Images {

		//fmt.Printf("Key---> %d \n", key)
		qsImageUrl := new(domain.QuestionImages)
		qsImageUrl.Imageurl = imageurl.Imageurl
		imagesArray = append(imagesArray, *qsImageUrl)

		// fmt.Printf("ImageURL: %s  \n", qsImageUrl.Imageurl)
		// fmt.Printf("%v \n --------------------", imagesArray)

	}

	singleQuestion.ID = _id
	singleQuestion.Question = questionText
	singleQuestion.Marks = marks
	singleQuestion.Answer = answer
	singleQuestion.Topics = topicsArray
	singleQuestion.Images = imagesArray

	fmt.Printf("Single Qs ->>> %v \t \n", singleQuestion)
	return qsUC.questionRepo.AddSingleTheoryQuestion(ctx, singleQuestion, metadataID)
	//return &singleQuestion, nil

}

func (qsUC *questionUsecase) GetMetadataById(requestCtx context.Context, username string, useremail string) ([]domain.MetadataBson, error) {

	ctx, cancel := context.WithTimeout(requestCtx, qsUC.contextTimeout)
	defer cancel()

	return qsUC.questionRepo.GetMetadataById(ctx, username, useremail)
}

func (qsUC *questionUsecase) UpdateMetadataById(requestCtx context.Context, receivedMetadata domain.MetadataBson, docID string) (int64, error) {

	ctx, cancel := context.WithTimeout(requestCtx, qsUC.contextTimeout)
	defer cancel()

	return qsUC.questionRepo.UpdateMetadataById(ctx, receivedMetadata, docID)
}

func (qsUC *questionUsecase) DeleteMetadataById(requestCtx context.Context, docID string) (int64, error) {

	ctx, cancel := context.WithTimeout(requestCtx, qsUC.contextTimeout)
	defer cancel()

	return qsUC.questionRepo.DeleteMetadataById(ctx, docID)
}

func (qsUC *questionUsecase) GetMCQsByMetadataID(requestCtx context.Context, metadataID string) ([]domain.DisplayQuestion, error) {

	ctx, cancel := context.WithTimeout(requestCtx, qsUC.contextTimeout)
	defer cancel()

	return qsUC.questionRepo.GetMCQsByMetadataID(ctx, metadataID)	
}

func (qsUC *questionUsecase) GetQuestion(requestCtx context.Context, objectID string) (domain.Question, error) {

	ctx, cancel := context.WithTimeout(requestCtx, qsUC.contextTimeout)
	defer cancel()

	return qsUC.questionRepo.GetQuestion(ctx, objectID)	
}

func (qsUC *questionUsecase) GetTheoryQuestion(requestCtx context.Context, objectID string) (domain.TheoryQuestion, error) {

	ctx, cancel := context.WithTimeout(requestCtx, qsUC.contextTimeout)
	defer cancel()

	return qsUC.questionRepo.GetTheoryQuestion(ctx, objectID)	
}

func (qsUC *questionUsecase) UpdateQuestion(requestCtx context.Context, updatedQuestion domain.Question, questionID string) (int64, error) {

	ctx, cancel := context.WithTimeout(requestCtx, qsUC.contextTimeout)
	defer cancel()

	return qsUC.questionRepo.UpdateQuestion(ctx, updatedQuestion, questionID)
}

func (qsUC *questionUsecase) UpdateTheoryQuestion(requestCtx context.Context, updatedQuestion domain.TheoryQuestion, questionID string) (int64, error) {

	ctx, cancel := context.WithTimeout(requestCtx, qsUC.contextTimeout)
	defer cancel()

	return qsUC.questionRepo.UpdateTheoryQuestion(ctx, updatedQuestion, questionID)
}

func (qsUC *questionUsecase) DeleteQuestion(requestCtx context.Context, metaID string, questionID string) (int64, error) {

	ctx, cancel := context.WithTimeout(requestCtx, qsUC.contextTimeout)
	defer cancel()

	return qsUC.questionRepo.DeleteQuestion(ctx, metaID, questionID)
}