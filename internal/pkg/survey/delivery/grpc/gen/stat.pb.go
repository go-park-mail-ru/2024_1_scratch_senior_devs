// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v5.26.1
// source: stat.proto

package gen

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type VoteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	QuestionId string `protobuf:"bytes,1,opt,name=QuestionId,proto3" json:"QuestionId,omitempty"`
	Vote       int32  `protobuf:"varint,2,opt,name=Vote,proto3" json:"Vote,omitempty"`
}

func (x *VoteRequest) Reset() {
	*x = VoteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stat_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VoteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VoteRequest) ProtoMessage() {}

func (x *VoteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_stat_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VoteRequest.ProtoReflect.Descriptor instead.
func (*VoteRequest) Descriptor() ([]byte, []int) {
	return file_stat_proto_rawDescGZIP(), []int{0}
}

func (x *VoteRequest) GetQuestionId() string {
	if x != nil {
		return x.QuestionId
	}
	return ""
}

func (x *VoteRequest) GetVote() int32 {
	if x != nil {
		return x.Vote
	}
	return 0
}

type GetSurveyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetSurveyRequest) Reset() {
	*x = GetSurveyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stat_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetSurveyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSurveyRequest) ProtoMessage() {}

func (x *GetSurveyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_stat_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSurveyRequest.ProtoReflect.Descriptor instead.
func (*GetSurveyRequest) Descriptor() ([]byte, []int) {
	return file_stat_proto_rawDescGZIP(), []int{1}
}

type VoteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *VoteResponse) Reset() {
	*x = VoteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stat_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VoteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VoteResponse) ProtoMessage() {}

func (x *VoteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_stat_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VoteResponse.ProtoReflect.Descriptor instead.
func (*VoteResponse) Descriptor() ([]byte, []int) {
	return file_stat_proto_rawDescGZIP(), []int{2}
}

type Question struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           string `protobuf:"bytes,1,opt,name=Id,proto3" json:"Id,omitempty"`
	Title        string `protobuf:"bytes,2,opt,name=Title,proto3" json:"Title,omitempty"`
	QuestionType string `protobuf:"bytes,3,opt,name=QuestionType,proto3" json:"QuestionType,omitempty"`
	Number       int64  `protobuf:"varint,4,opt,name=number,proto3" json:"number,omitempty"`
	SurveyId     string `protobuf:"bytes,5,opt,name=SurveyId,proto3" json:"SurveyId,omitempty"`
}

func (x *Question) Reset() {
	*x = Question{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stat_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Question) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Question) ProtoMessage() {}

func (x *Question) ProtoReflect() protoreflect.Message {
	mi := &file_stat_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Question.ProtoReflect.Descriptor instead.
func (*Question) Descriptor() ([]byte, []int) {
	return file_stat_proto_rawDescGZIP(), []int{3}
}

func (x *Question) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Question) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Question) GetQuestionType() string {
	if x != nil {
		return x.QuestionType
	}
	return ""
}

func (x *Question) GetNumber() int64 {
	if x != nil {
		return x.Number
	}
	return 0
}

func (x *Question) GetSurveyId() string {
	if x != nil {
		return x.SurveyId
	}
	return ""
}

type GetSurveyResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Questions []*Question `protobuf:"bytes,1,rep,name=Questions,proto3" json:"Questions,omitempty"`
}

func (x *GetSurveyResponse) Reset() {
	*x = GetSurveyResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stat_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetSurveyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSurveyResponse) ProtoMessage() {}

func (x *GetSurveyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_stat_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSurveyResponse.ProtoReflect.Descriptor instead.
func (*GetSurveyResponse) Descriptor() ([]byte, []int) {
	return file_stat_proto_rawDescGZIP(), []int{4}
}

func (x *GetSurveyResponse) GetQuestions() []*Question {
	if x != nil {
		return x.Questions
	}
	return nil
}

type CreateSurveyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Questions []*Question `protobuf:"bytes,1,rep,name=Questions,proto3" json:"Questions,omitempty"`
}

func (x *CreateSurveyRequest) Reset() {
	*x = CreateSurveyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stat_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateSurveyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateSurveyRequest) ProtoMessage() {}

func (x *CreateSurveyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_stat_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateSurveyRequest.ProtoReflect.Descriptor instead.
func (*CreateSurveyRequest) Descriptor() ([]byte, []int) {
	return file_stat_proto_rawDescGZIP(), []int{5}
}

func (x *CreateSurveyRequest) GetQuestions() []*Question {
	if x != nil {
		return x.Questions
	}
	return nil
}

type CreateSurveyResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CreateSurveyResponse) Reset() {
	*x = CreateSurveyResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stat_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateSurveyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateSurveyResponse) ProtoMessage() {}

func (x *CreateSurveyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_stat_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateSurveyResponse.ProtoReflect.Descriptor instead.
func (*CreateSurveyResponse) Descriptor() ([]byte, []int) {
	return file_stat_proto_rawDescGZIP(), []int{6}
}

var File_stat_proto protoreflect.FileDescriptor

var file_stat_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x6e, 0x6f,
	0x74, 0x65, 0x22, 0x41, 0x0a, 0x0b, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x49,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x56, 0x6f, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x04, 0x56, 0x6f, 0x74, 0x65, 0x22, 0x12, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x53, 0x75, 0x72, 0x76,
	0x65, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x0e, 0x0a, 0x0c, 0x56, 0x6f, 0x74,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x88, 0x01, 0x0a, 0x08, 0x51, 0x75,
	0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x22, 0x0a, 0x0c,
	0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0c, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x53, 0x75, 0x72, 0x76,
	0x65, 0x79, 0x49, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x53, 0x75, 0x72, 0x76,
	0x65, 0x79, 0x49, 0x64, 0x22, 0x41, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x53, 0x75, 0x72, 0x76, 0x65,
	0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2c, 0x0a, 0x09, 0x51, 0x75, 0x65,
	0x73, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x6e,
	0x6f, 0x74, 0x65, 0x2e, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x09, 0x51, 0x75,
	0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x43, 0x0a, 0x13, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x53, 0x75, 0x72, 0x76, 0x65, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2c,
	0x0a, 0x09, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x0e, 0x2e, 0x6e, 0x6f, 0x74, 0x65, 0x2e, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x09, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x16, 0x0a, 0x14,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x75, 0x72, 0x76, 0x65, 0x79, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x32, 0xc0, 0x01, 0x0a, 0x04, 0x53, 0x74, 0x61, 0x74, 0x12, 0x3e, 0x0a,
	0x09, 0x47, 0x65, 0x74, 0x53, 0x75, 0x72, 0x76, 0x65, 0x79, 0x12, 0x16, 0x2e, 0x6e, 0x6f, 0x74,
	0x65, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x75, 0x72, 0x76, 0x65, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x17, 0x2e, 0x6e, 0x6f, 0x74, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x75, 0x72,
	0x76, 0x65, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x2f, 0x0a,
	0x04, 0x56, 0x6f, 0x74, 0x65, 0x12, 0x11, 0x2e, 0x6e, 0x6f, 0x74, 0x65, 0x2e, 0x56, 0x6f, 0x74,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x6e, 0x6f, 0x74, 0x65, 0x2e,
	0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x47,
	0x0a, 0x0c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x75, 0x72, 0x76, 0x65, 0x79, 0x12, 0x19,
	0x2e, 0x6e, 0x6f, 0x74, 0x65, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x75, 0x72, 0x76,
	0x65, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x6e, 0x6f, 0x74, 0x65,
	0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x75, 0x72, 0x76, 0x65, 0x79, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x2e, 0x5a, 0x2c, 0x2e, 0x2f, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x73, 0x75, 0x72, 0x76, 0x65, 0x79,
	0x2f, 0x64, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x79, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x67,
	0x65, 0x6e, 0x2f, 0x3b, 0x67, 0x65, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_stat_proto_rawDescOnce sync.Once
	file_stat_proto_rawDescData = file_stat_proto_rawDesc
)

func file_stat_proto_rawDescGZIP() []byte {
	file_stat_proto_rawDescOnce.Do(func() {
		file_stat_proto_rawDescData = protoimpl.X.CompressGZIP(file_stat_proto_rawDescData)
	})
	return file_stat_proto_rawDescData
}

var file_stat_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_stat_proto_goTypes = []interface{}{
	(*VoteRequest)(nil),          // 0: note.VoteRequest
	(*GetSurveyRequest)(nil),     // 1: note.GetSurveyRequest
	(*VoteResponse)(nil),         // 2: note.VoteResponse
	(*Question)(nil),             // 3: note.Question
	(*GetSurveyResponse)(nil),    // 4: note.GetSurveyResponse
	(*CreateSurveyRequest)(nil),  // 5: note.CreateSurveyRequest
	(*CreateSurveyResponse)(nil), // 6: note.CreateSurveyResponse
}
var file_stat_proto_depIdxs = []int32{
	3, // 0: note.GetSurveyResponse.Questions:type_name -> note.Question
	3, // 1: note.CreateSurveyRequest.Questions:type_name -> note.Question
	1, // 2: note.Stat.GetSurvey:input_type -> note.GetSurveyRequest
	0, // 3: note.Stat.Vote:input_type -> note.VoteRequest
	5, // 4: note.Stat.CreateSurvey:input_type -> note.CreateSurveyRequest
	4, // 5: note.Stat.GetSurvey:output_type -> note.GetSurveyResponse
	2, // 6: note.Stat.Vote:output_type -> note.VoteResponse
	6, // 7: note.Stat.CreateSurvey:output_type -> note.CreateSurveyResponse
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_stat_proto_init() }
func file_stat_proto_init() {
	if File_stat_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_stat_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VoteRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_stat_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetSurveyRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_stat_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VoteResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_stat_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Question); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_stat_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetSurveyResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_stat_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateSurveyRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_stat_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateSurveyResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_stat_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_stat_proto_goTypes,
		DependencyIndexes: file_stat_proto_depIdxs,
		MessageInfos:      file_stat_proto_msgTypes,
	}.Build()
	File_stat_proto = out.File
	file_stat_proto_rawDesc = nil
	file_stat_proto_goTypes = nil
	file_stat_proto_depIdxs = nil
}
