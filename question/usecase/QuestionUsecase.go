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

func (qsUC *questionUsecase) Save(c context.Context, allMcqs *domain.MCQModel, username string, useremail string) (err error) {
	ctx, cancel := context.WithTimeout(c, qsUC.contextTimeout)
	defer cancel()

	metadataBson := new(domain.MetadataBson)
	questionSet := []interface{}{}
	questionHexIds := make([]string, 0)

	for key, allQuestions := range *allMcqs {

		singleQuestion := new(domain.Question)

		fmt.Printf("\n *************************** \n")
		if key == 0 {

			metadataBson.System = allQuestions.System
			metadataBson.Board = allQuestions.Board
			metadataBson.Subject = allQuestions.Subject
			metadataBson.Year = allQuestions.Year
			metadataBson.Month = allQuestions.Month
			metadataBson.Series = allQuestions.Series
			metadataBson.Paper = allQuestions.Paper
			metadataBson.Username = username
			metadataBson.Useremail = useremail

		} else {
			_id := primitive.NewObjectID()
			questionText := allQuestions.Questions
			marks := allQuestions.Marks
			optionsArray := make([]domain.QuestionOptions, 0)
			topicsArray := make([]domain.QuestionTopic, 0)

			// fmt.Printf("Question:-->  %s --- Marks:--> %s \n ", questionText, marks)
			for _, option := range allQuestions.Options {

				//fmt.Printf("Key---> %d \n", key)
				qsOption := new(domain.QuestionOptions)
				qsOption.Option = option.Option
				qsOption.Correct = option.Correct
				optionsArray = append(optionsArray, *qsOption)

				// fmt.Printf("Option: %s --- Correct: %t \n", qsOption.Option, qsOption.Correct)
				// fmt.Printf("%v \n --------------------", optionsArray)

			}

			for _, topic := range allQuestions.Topic {

				//fmt.Printf("Key---> %d \n", key)
				qsTopic := new(domain.QuestionTopic)
				qsTopic.Topic = topic.Topic
				topicsArray = append(topicsArray, *qsTopic)

				// fmt.Printf("Topic: %s  \n", qsTopic.Topic)
				// fmt.Printf("%v \n --------------------", topicsArray)

			}

			singleQuestion.ID = _id
			singleQuestion.Questions = questionText
			singleQuestion.Marks = marks
			singleQuestion.Options = optionsArray
			singleQuestion.Topics = topicsArray

			questionSet = append(questionSet, singleQuestion)

			//adding hexID to metadata
			questionHexIds = append(questionHexIds, _id.Hex())

			fmt.Printf("Single Qs ->>> %v \n", singleQuestion)
		}
	}

	metadataBson.QuestionHexIds = questionHexIds

	fmt.Println("*******************************________________**********************************")
	fmt.Printf("\n %#v", questionSet)
	fmt.Printf("\n ID---- %#v", questionHexIds)

	qsUC.questionRepo.SaveQuestionMetadata(ctx, metadataBson)
	qsUC.questionRepo.SaveAllQuestions(ctx, questionSet)
	return
}
