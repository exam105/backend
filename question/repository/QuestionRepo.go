package repository

import (
	"github.com/exam105-UPD/backend/logging"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"context"
	"fmt"
	"log"

	"github.com/exam105-UPD/backend/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type questionRepo struct {
	Conn *mongo.Client
}

// This will create an object that represent the question.Repository interface
func NewQuestionRepository(Conn *mongo.Client) domain.QuestionRepository {
	return &questionRepo{Conn}
}

func (db *questionRepo) SaveQuestionMetadata(ctx context.Context, qsMetaData *domain.MetadataBson) (error){

	database := db.Conn.Database("exam105")
	metadataCollection := database.Collection("metadata")
	insertResult, err := metadataCollection.InsertOne(ctx, qsMetaData)
	if err != nil {
		log.Println( logging.MSG_InsertUnsuccessful, err.Error())
		return fmt.Errorf(logging.MSG_InsertUnsuccessful + "\n ID: %s \t" + err.Error(), insertResult.InsertedID)
	}

	fmt.Println("Inserted a single metadata document: ", insertResult)
	return nil
}

func (db *questionRepo) SaveAllQuestions(ctx context.Context, questions []interface{}) (int64, error){

	database := db.Conn.Database("exam105")
	questionsCollection := database.Collection("questions")

	// create the slice of write models
	var writes []mongo.WriteModel
	for _, ins := range questions {
		model := mongo.NewInsertOneModel().SetDocument(ins)
		writes = append(writes, model)
	}

	// run bulk write
	res, err := questionsCollection.BulkWrite(ctx, writes)
	if err != nil {
		log.Println( logging.MSG_BulkwriteFailed, err.Error())
		return -1, fmt.Errorf(logging.MSG_BulkwriteFailed + "\n ID: %s \t" + err.Error(), res.InsertedCount)
	}

	fmt.Printf(
		"insert: %d, updated: %d, deleted: %d",
		res.InsertedCount,
		res.ModifiedCount,
		res.DeletedCount,
	)

	//insertResult, err := questionsCollection.InsertMany(ctx, questions)

	return res.InsertedCount,nil
}

func (db *questionRepo) GetMetadataById(ctx context.Context, username string, useremail string) ([]domain.MetadataBson, error) {

	database := db.Conn.Database("exam105")
	metadataCollection := database.Collection("metadata")

	cursor, err := metadataCollection.Find(ctx,
		bson.D{
			{"username", username},
			{"useremail", useremail},
		})

	if err != nil {
		log.Println( logging.MSG_DocumentNotFound, err.Error())
		return []domain.MetadataBson{}, fmt.Errorf(logging.MSG_DocumentNotFound + "\n ID: %s \t" + err.Error(), cursor.ID)
	}
	defer cursor.Close(ctx)

	var metadata []domain.MetadataBson
	if err = cursor.All(ctx, &metadata); err != nil {
		log.Println( logging.MSG_MappingFailure, err.Error())
		return []domain.MetadataBson{}, fmt.Errorf(logging.MSG_MappingFailure + "\n ID: %s \t" + err.Error(), metadata)
	}

	if len(metadata) == 0 {
		log.Println( logging.MSG_EmptyMetadata, err.Error())
		return []domain.MetadataBson{}, fmt.Errorf(logging.MSG_EmptyMetadata + "\n ID: %s \t" + err.Error(), metadata)
	}

	return metadata, nil
}

func (db *questionRepo) UpdateMetadataById(ctx context.Context, receivedMetadata domain.MetadataBson, docID string) (int64, error) {

	database := db.Conn.Database("exam105")
	metadataCollection := database.Collection("metadata")

	mongoID, err := primitive.ObjectIDFromHex(docID)
	
	if err != nil {
		log.Println( logging.MSG_ConversionUnsuccessful, err.Error())
		return -1, fmt.Errorf(logging.MSG_ConversionUnsuccessful + "\n ID: %s \t" + err.Error(), docID)
	} 

	result, err := metadataCollection.UpdateOne(
		ctx,
		bson.M{"_id": mongoID},
		bson.D{
			{"$set", bson.D{
				{Key: "subject", Value: receivedMetadata.Subject},
				{Key: "system", Value: receivedMetadata.System},
				{Key: "board", Value: receivedMetadata.Board},
				{Key: "series", Value: receivedMetadata.Series},
				{Key: "paper", Value: receivedMetadata.Paper},
				{Key: "year", Value: receivedMetadata.Year},
				{Key: "month", Value: receivedMetadata.Month},
			}},
		},
	)
	if err != nil {
		log.Println( logging.MSG_UpdateUnsuccessful, err.Error())
		return -1, fmt.Errorf(logging.MSG_UpdateUnsuccessful + "\n ID: %s \t" + err.Error(), mongoID)
		
	}

	fmt.Printf(
		"insert: %d, updated: %d, deleted: %d /n",
		result.MatchedCount,
		result.ModifiedCount,
		result.UpsertedCount,
	)
	return result.ModifiedCount, nil
}

func (db *questionRepo) DeleteMetadataById(ctx context.Context, docID string) (int64, error) {

	var questionHexIDCollection []primitive.ObjectID
	database := db.Conn.Database("exam105")
	metadataCollection := database.Collection("metadata")
	questionCollection := database.Collection("questions")

	metadataID, err := primitive.ObjectIDFromHex(docID)

	if err != nil {
		log.Println( logging.MSG_ConversionUnsuccessful, err.Error())
		return -1, fmt.Errorf(logging.MSG_ConversionUnsuccessful + "\n ID: %s \t" + err.Error(), docID)
	} 

	// 1. First all the Questions based on HexIDs - Metadata (Exam Paper)
	var metadataSingleRecord domain.MetadataBson
	err = metadataCollection.FindOne(ctx, 
		bson.M{"_id": metadataID}).Decode(&metadataSingleRecord)
	
	if err != nil {
		log.Println( logging.MSG_WrongDocumentID, err.Error())
		return -1, fmt.Errorf(logging.MSG_WrongDocumentID + "\n ID: %s \t" + err.Error(), metadataID)
	}

	//Collecting all the questions IDs of the paper
	for k, _ := range metadataSingleRecord.QuestionHexIds {
	
		fmt.Printf("Question HEX ID: %s\t\n", metadataSingleRecord.QuestionHexIds[k])
		questionHexID := metadataSingleRecord.QuestionHexIds[k]
		questionID, err := primitive.ObjectIDFromHex(questionHexID)
		if err != nil {
			log.Println( logging.MSG_ConversionUnsuccessful, err.Error())
			return -1, fmt.Errorf(logging.MSG_ConversionUnsuccessful + "\n ID: %s \t" + err.Error(), questionID)
		}

		//Populating a slice of all the questionIDs
		questionHexIDCollection = append(questionHexIDCollection, questionID)
	}

	filter := bson.D{
		{"_id", bson.D{{"$in", questionHexIDCollection}}}, }
		
	deleteResult, err := questionCollection.DeleteMany(ctx, filter)
	if err != nil {
		log.Println( logging.MSG_DeleteUnsuccessful, err.Error())
		return -1, fmt.Errorf(logging.MSG_DeleteUnsuccessful + "\n ID: %s \t" + err.Error(), deleteResult.DeletedCount)
	}

	fmt.Println("All question of Metadata Deleted: ", deleteResult)

	// 2. Delete the Metadata from Metadata collection
	result, err := metadataCollection.DeleteOne(ctx, bson.M{"_id": metadataID})

	if err != nil {
		log.Println( logging.MSG_DeleteUnsuccessful, err.Error())
		return -1, fmt.Errorf(logging.MSG_DeleteUnsuccessful + "\n ID: %s \t" + err.Error(), metadataID)
	} 

	return deleteResult.DeletedCount + result.DeletedCount, nil
}

func (db *questionRepo) GetMCQsByMetadataID(ctx context.Context, docID string) ([]domain.DisplayQuestion, error) {

	database := db.Conn.Database("exam105")
	metadataCollection := database.Collection("metadata")
	questionsCollection := database.Collection("questions")
	metadataID, err := primitive.ObjectIDFromHex(docID)

	//Metadata (Exam Paper)
	var metadataSingleRecord domain.MetadataBson
	err = metadataCollection.FindOne(ctx, 
		bson.M{"_id": metadataID}).Decode(&metadataSingleRecord)
	
	if err != nil {
		log.Println( logging.MSG_WrongDocumentID, err )
		return nil, fmt.Errorf(logging.MSG_WrongDocumentID + "\n ID: %s \t" + err.Error(), docID)
	}
	
	fmt.Printf("Host   : %s\t\n", metadataSingleRecord.Subject)

	//Selecting all the questions of the paper
	questionToDisplay := []domain.DisplayQuestion{}

	for k, _ := range metadataSingleRecord.QuestionHexIds {
	
		fmt.Printf("Question HEX ID: %s\t\n", metadataSingleRecord.QuestionHexIds[k])	
		questionHexID := metadataSingleRecord.QuestionHexIds[k]
		questionID, err := primitive.ObjectIDFromHex(questionHexID)

		var question domain.DisplayQuestion
		err = questionsCollection.FindOne(ctx, 
			bson.M{"_id": questionID}).Decode(&question)
				
		if err != nil {
			log.Println( logging.MSG_WrongDocumentID, err)
			return nil, fmt.Errorf(logging.MSG_WrongDocumentID + "\n ID: %s \t" + err.Error(), questionID)
		}
	
		questionToDisplay = append(questionToDisplay, question)
	}

	fmt.Printf("All questions: %+v\n", questionToDisplay)
	
	return questionToDisplay, nil

}

func (db *questionRepo) GetQuestion(ctx context.Context, questionID string) (domain.Question, error) {

	database := db.Conn.Database("exam105")
	questionsCollection := database.Collection("questions")
	question, err := primitive.ObjectIDFromHex(questionID)

	var singleQuestion domain.Question
	err = questionsCollection.FindOne(ctx, 
		bson.M{"_id": question}).Decode(&singleQuestion)

	if err != nil {
		log.Println( logging.MSG_WrongDocumentID, err.Error())
		return domain.Question{}, fmt.Errorf(logging.MSG_WrongDocumentID + "\n ID: %s \t" + err.Error(), questionID)
	}

	fmt.Printf("Questions: %+v \n", singleQuestion)
	return singleQuestion, nil
}

func (db *questionRepo) GetTheoryQuestion(ctx context.Context, questionID string) (domain.TheoryQuestion, error) {

	database := db.Conn.Database("exam105")
	questionsCollection := database.Collection("questions")
	question, err := primitive.ObjectIDFromHex(questionID)

	var singleQuestion domain.TheoryQuestion
	err = questionsCollection.FindOne(ctx, 
		bson.M{"_id": question}).Decode(&singleQuestion)

	if err != nil {
		log.Println( logging.MSG_WrongDocumentID, err.Error())
		return domain.TheoryQuestion{}, fmt.Errorf(logging.MSG_WrongDocumentID + "\n ID: %s \t" + err.Error(), questionID)
	}

	fmt.Printf("Questions: %+v \n", singleQuestion)
	return singleQuestion, nil
}

func (db *questionRepo) UpdateQuestion(ctx context.Context, updatedQuestion domain.Question, questionID string) (int64, error) {

	database := db.Conn.Database("exam105")
	questionCollection := database.Collection("questions")

	questionId, err := primitive.ObjectIDFromHex(questionID)

	if err != nil {
		log.Println( logging.MSG_ConversionUnsuccessful, err.Error())
		return -1, fmt.Errorf(logging.MSG_ConversionUnsuccessful + "\n ID: %s \t" + err.Error(), questionID)
	} 

	// Use it's ID to replace
	filter := bson.M{"_id": questionId}
	result, err := questionCollection.ReplaceOne(ctx, filter, updatedQuestion)

	if err != nil {
		log.Println( logging.MSG_UpdateUnsuccessful, err.Error())
		return -1, fmt.Errorf(logging.MSG_UpdateUnsuccessful + "\n ID: %s \t" + err.Error(), result)
	}

	fmt.Printf(
		"insert: %d, updated: %d, deleted: %d /n",
		result.MatchedCount,
		result.ModifiedCount,
		result.UpsertedCount,
	)
	return result.ModifiedCount, nil
}

func (db *questionRepo) UpdateTheoryQuestion(ctx context.Context, updatedQuestion domain.TheoryQuestion, questionID string) (int64, error) {

	database := db.Conn.Database("exam105")
	questionCollection := database.Collection("questions")

	questionId, err := primitive.ObjectIDFromHex(questionID)

	if err != nil {
		log.Println( logging.MSG_ConversionUnsuccessful, err.Error())
		return -1, fmt.Errorf(logging.MSG_ConversionUnsuccessful + "\n ID: %s \t" + err.Error(), questionID)
	} 

	// Use it's ID to replace
	filter := bson.M{"_id": questionId}
	result, err := questionCollection.ReplaceOne(ctx, filter, updatedQuestion)

	if err != nil {
		log.Println( logging.MSG_UpdateUnsuccessful, err.Error())
		return -1, fmt.Errorf(logging.MSG_UpdateUnsuccessful + "\n ID: %s \t" + err.Error(), result)
	}

	fmt.Printf(
		"insert: %d, updated: %d, deleted: %d /n",
		result.MatchedCount,
		result.ModifiedCount,
		result.UpsertedCount,
	)
	return result.ModifiedCount, nil
}

func (db *questionRepo) DeleteQuestion(ctx context.Context, metaID string, questionID string) (int64, error) {

	database := db.Conn.Database("exam105")
	questionCollection := database.Collection("questions")
	metadataCollection := database.Collection("metadata")

	questionHexID, ques_err := primitive.ObjectIDFromHex(questionID)
	metaHexID, meta_err := primitive.ObjectIDFromHex(metaID)

	if ques_err != nil {
		log.Println(logging.MSG_WrongDocumentID, ques_err.Error())
		return -1, fmt.Errorf(logging.MSG_WrongDocumentID + "\n ID: %s \t" + ques_err.Error(), questionID)
	} else if meta_err != nil {
		log.Println(logging.MSG_WrongDocumentID, meta_err.Error() )
		return -1, fmt.Errorf(logging.MSG_WrongDocumentID + "\n ID: %s \t" + meta_err.Error(), metaID)
	}
	
	// Start Transaction
	// var session mongo.Session
	// if session, err = db.Conn.StartSession(); err != nil {
    //     log.Fatal(err)
    // }
    // if err = session.StartTransaction(); err != nil {
    //     log.Fatal(err)
    // }

	// 1. First we need to delete the QuestionID from the array in Metadata collection
	var transactionStatus int64
	var metadataSingleRecord domain.MetadataBson
	err := metadataCollection.FindOne(ctx, bson.M{"_id": metaHexID}).Decode(&metadataSingleRecord)
	
	if err != nil {
		log.Println(logging.MSG_DocumentNotFound, err.Error() )
		return -1, fmt.Errorf(logging.MSG_DocumentNotFound + "\n ID: %s \t" + err.Error(), metaHexID)
	}
	
	fmt.Printf("Full list   : %s\t \n", metadataSingleRecord.QuestionHexIds)
	var questionHexs = make([]string,0)
	questionHexs = append(questionHexs, metadataSingleRecord.QuestionHexIds...)

	for k, value := range questionHexs{
		if  value == questionID {

			//Removing QuestionHexID from Metadata Collection
			questionHexs = removeIndex(questionHexs, k)
			break
		} 
	}
	// 1.1 Updating Metadata collection after removing ID
	filter := bson.M{"_id": bson.M{"$eq": metaHexID}}
	update := bson.M{"$set": bson.M{"question_hex_ids": questionHexs}}
	updated, updateErr := metadataCollection.UpdateOne(ctx, filter, update)
	
	if updateErr != nil {
		log.Println(logging.MSG_DeleteUnsuccessful, updateErr.Error() )
		return -1, fmt.Errorf(logging.MSG_DeleteUnsuccessful + "\n ID: %s \t" + updateErr.Error(), questionHexs)
	}

	transactionStatus = updated.ModifiedCount
	fmt.Printf("Question_HEX_IDs   : %s \t %d \n", questionHexs, len(questionHexs))

	// 2. Delete the Question Document from Question Collection
	deleted, deleteErr := questionCollection.DeleteOne(ctx, bson.M{"_id": questionHexID})

	if deleteErr != nil {
		log.Println(logging.MSG_DeleteUnsuccessful, deleteErr.Error() )
		return -1, fmt.Errorf(logging.MSG_DeleteUnsuccessful + "\n ID: %s \t" + deleteErr.Error(), questionHexID)
	} 

	transactionStatus = transactionStatus + deleted.DeletedCount

	//**************Transaction Needs Replica-Set**********************
 	
	// var transactionStatus int64
	// if err = mongo.WithSession(ctx, session, func(mongoSession mongo.SessionContext) error {

	// 	// 1.1 Updating Metadata collection after removing ID
	// 	filter := bson.M{"_id": bson.M{"$eq": metaHexID}}
	// 	update := bson.M{"$set": bson.M{"question_hex_ids": questionHexs}}
	// 	updated, updateErr := metadataCollection.UpdateOne(mongoSession, filter, update)
		
	// 	if updateErr != nil {
	// 		log.Fatal("Delete Metadata ERROR: ", updateErr)
	// 		mongoSession.AbortTransaction(mongoSession)
	// 	} else {

	// 		if updated.ModifiedCount != 1 {
	// 			fmt.Println("Delete failed. Expected 1 but got ", updated)
	// 		} else {
	// 			fmt.Println("Deleteed: ", updated)
	// 		}
	// 	}	
	// 	transactionStatus = updated.ModifiedCount
	// 	fmt.Printf("Question_HEX_IDs   : %s \t %d \n", questionHexs, len(questionHexs))

	// 	// 2. Delete the Question Document from Question Collection
	// 	deleted, deleteErr := questionCollection.DeleteOne(mongoSession, bson.M{"_id": questionHexID})

	// 	if deleteErr != nil {
	// 		log.Fatal("Delete Question ERROR:", deleteErr)
	// 		mongoSession.AbortTransaction(mongoSession)
	// 	} else {

	// 		if deleted.DeletedCount == 0 {
	// 			fmt.Println("Delete document not found:", deleted)
	// 		} else {
	// 			fmt.Println("Delete Result:", deleted)
	// 		}
	// 	}
	// 	transactionStatus = transactionStatus + deleted.DeletedCount
	// 	if err = session.CommitTransaction(mongoSession); err != nil {
    //         log.Fatal(err)
    //     }
    //     return nil
    // }); err != nil {
    //     log.Fatal(err)
    // }
    // session.EndSession(ctx) // End transaction

	return transactionStatus, nil
}

func (db *questionRepo) AddSingleQuestion(ctx context.Context, singleQuestion *domain.Question, metadataID string) (int64, error) {

 	database := db.Conn.Database("exam105")
	questionsCollection := database.Collection("questions")
	metadataCollection := database.Collection("metadata")

	metaHexID, meta_err := primitive.ObjectIDFromHex(metadataID)

	if meta_err != nil {
		log.Println(logging.MSG_ConversionUnsuccessful, meta_err.Error() )
		return -1, fmt.Errorf(logging.MSG_ConversionUnsuccessful + "\n ID: %s \t" + meta_err.Error(), metadataID)		
	}

	// 1. Adding questionID to the Metadata question array
	var metadataSingleRecord domain.MetadataBson
	err := metadataCollection.FindOne(ctx, bson.M{"_id": metaHexID}).Decode(&metadataSingleRecord)
	
	if err != nil {
		log.Println( logging.MSG_DocumentNotFound, err.Error())
		return -1, fmt.Errorf(logging.MSG_DocumentNotFound + "\n ID: %s \t" + err.Error(), metaHexID)			
	}
	
	// fmt.Printf("Full list   : %s\t \n", metadataSingleRecord.QuestionHexIds)

	var questionHexs = make([]string,0)
	questionHexs = append(questionHexs, metadataSingleRecord.QuestionHexIds...)
	questionHexs = append(questionHexs, singleQuestion.ID.Hex())	//New QuestionHexID added to the slice

	fmt.Println("Question IDs: ", questionHexs)	

	// 1. Updating Metadata collection to add new QuestionID
	filter := bson.M{"_id": bson.M{"$eq": metaHexID}}
	questionsSlice := bson.M{"$set": bson.M{"question_hex_ids": questionHexs}}
	updated, updateErr := metadataCollection.UpdateOne(ctx, filter, questionsSlice)	

	if updateErr != nil {
		log.Println( logging.MSG_UpdateUnsuccessful, updateErr.Error())
		return -1, fmt.Errorf(logging.MSG_UpdateUnsuccessful + "\n ID: %s \t" + updateErr.Error(), updated.ModifiedCount)		
	}
	
	// 2.Adding question to the Question Collection
	insert, insertErr := questionsCollection.InsertOne(ctx, singleQuestion)
	if insertErr != nil {
		log.Println( logging.MSG_InsertUnsuccessful, insertErr.Error())
		return -1, fmt.Errorf(logging.MSG_InsertUnsuccessful + "\n ID: %s \t" + updateErr.Error(), insert.InsertedID)		
	} 

	fmt.Printf(
		"insert: %d, updated: %d",
		insert.InsertedID,
		updated.ModifiedCount,
	)

	return  updated.ModifiedCount, nil
}

func removeIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}