package logging

import (
	"github.com/spf13/viper"
)


var MSG_ConversionUnsuccessful string
var MSG_DocumentNotFound string
var MSG_UpdateUnsuccessful string
var MSG_UpdateSuccessful string
var MSG_WrongDocumentID string
var MSG_DeleteUnsuccessful string
var MSG_DeleteSuccessful string
var MSG_InsertUnsuccessful string
var MSG_InsertSuccessful string
var MSG_EmptyMetadata string
var MSG_MappingFailure string
var MSG_BulkwriteFailed string
var MSG_LoginFailed string
var MSG_DocumentFound string


func InitializeMessages(){
	
	MSG_ConversionUnsuccessful = viper.GetString(`errMessages.conversionUnsuccessful`)
	MSG_DocumentNotFound = viper.GetString(`errMessages.documentNotFound`)
	MSG_UpdateUnsuccessful = viper.GetString(`errMessages.updateUnsuccessful`)
	MSG_UpdateSuccessful = viper.GetString(`errMessages.updateSuccessful`)
	MSG_WrongDocumentID = viper.GetString(`errMessages.wrongDocumentID`)
	MSG_DeleteUnsuccessful = viper.GetString(`errMessages.deleteUnsuccessful`)
	MSG_DeleteSuccessful = viper.GetString(`errMessages.deleteSuccessful`)
	MSG_InsertUnsuccessful = viper.GetString(`errMessages.insertUnsuccessful`)
	MSG_InsertSuccessful = viper.GetString(`errMessages.insertSuccessful`)
	MSG_EmptyMetadata = viper.GetString(`errMessages.emptyMetadata`)
	MSG_MappingFailure = viper.GetString(`errMessages.mappingFailure`)
	MSG_LoginFailed = viper.GetString(`errMessages.loginFailed`)

	// Success Messages
	MSG_DocumentFound = viper.GetString(`successMessages.documentFound`)

}
