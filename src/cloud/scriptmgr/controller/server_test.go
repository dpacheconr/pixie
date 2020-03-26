package controller

import (
	"context"
	"fmt"
	"testing"

	"pixielabs.ai/pixielabs/src/cloud/scriptmgr/scriptmgrpb"

	mock_controller "pixielabs.ai/pixielabs/src/cloud/scriptmgr/controller/mock"
	statuspb "pixielabs.ai/pixielabs/src/common/base/proto"
	"pixielabs.ai/pixielabs/src/shared/scriptspb"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const visFuncsQuery = `
import px
@px.vis.vega("vega spec for f")
def f(start_time: px.Time, end_time: px.Time, svc: str):
  """Doc string for f"""
  return 1

@px.vis.vega("vega spec for g")
def g(a: int, b: float):
  """Doc string for g"""
	return 1
@px.vis.vega("vega spec for h")
def h(a: int, b: float):
	"""Doc string for h"""
	return 1
`

const visFuncsInfoPBStr = `
doc_string_map {
  key: "f"
  value: "Doc string for f"
}
doc_string_map {
  key: "g"
  value: "Doc string for g"
}
doc_string_map {
  key: "h"
  value: "Doc string for h"
}
vis_spec_map {
  key: "f"
  value {
    vega_spec: "vega spec for f"
  }
}
vis_spec_map {
  key: "g"
  value {
    vega_spec: "vega spec for g"
  }
}
vis_spec_map {
  key: "h"
  value {
    vega_spec: "vega spec for h"
  }
}
fn_args_map {
  key: "f"
  value {
    args {
      data_type: TIME64NS
      name: "start_time"
    }
    args {
      data_type: TIME64NS
      name: "end_time"
    }
    args {
      data_type: STRING
      name: "svc"
    }
  }
}
fn_args_map {
  key: "g"
  value {
    args {
      data_type: INT64
      name: "a"
    }
    args {
      data_type: FLOAT64
      name: "b"
    }
  }
}
fn_args_map {
  key: "h"
  value {
    args {
      data_type: INT64
      name: "a"
    }
    args {
      data_type: FLOAT64
      name: "b"
    }
  }
}
`

const expectedVisFuncsInfoSingleFilter = `
doc_string_map {
  key: "f"
  value: "Doc string for f"
}
vis_spec_map {
  key: "f"
  value {
    vega_spec: "vega spec for f"
  }
}
fn_args_map {
  key: "f"
  value {
    args {
      data_type: TIME64NS
      name: "start_time"
    }
    args {
      data_type: TIME64NS
      name: "end_time"
    }
    args {
      data_type: STRING
      name: "svc"
    }
  }
}
`

const expectedVisFuncsInfoMultiFilter = `
doc_string_map {
  key: "f"
  value: "Doc string for f"
}
doc_string_map {
  key: "g"
  value: "Doc string for g"
}
vis_spec_map {
  key: "f"
  value {
    vega_spec: "vega spec for f"
  }
}
vis_spec_map {
  key: "g"
  value {
    vega_spec: "vega spec for g"
  }
}
fn_args_map {
  key: "f"
  value {
    args {
      data_type: TIME64NS
      name: "start_time"
    }
    args {
      data_type: TIME64NS
      name: "end_time"
    }
    args {
      data_type: STRING
      name: "svc"
    }
  }
}
fn_args_map {
  key: "g"
  value {
    args {
      data_type: INT64
      name: "a"
    }
    args {
      data_type: FLOAT64
      name: "b"
    }
  }
}
`

func TestServerExtractVisFuncsInfo_NoFuncFilter(t *testing.T) {

	// Set up mocks.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	planner := mock_controller.NewMockPlanner(ctrl)

	visFuncsInfo := &scriptspb.VisFuncsInfo{}
	if err := proto.UnmarshalText(visFuncsInfoPBStr, visFuncsInfo); err != nil {
		t.Fatalf("Failed to unmarshal visfuncsinfo, err: %s", err)
	}
	expected := &scriptspb.VisFuncsInfo{}
	if err := proto.UnmarshalText(visFuncsInfoPBStr, expected); err != nil {
		t.Fatalf("Failed to unmarshal expected visfuncsinfo, err: %s", err)
	}

	planner.
		EXPECT().
		ExtractVisFuncsInfo(visFuncsQuery).
		Return(&scriptspb.VisFuncsInfoResult{
			Info:   visFuncsInfo,
			Status: &statuspb.Status{ErrCode: statuspb.OK},
		}, nil)

	req := &scriptmgrpb.ExtractVisFuncsInfoRequest{
		Script:    visFuncsQuery,
		FuncNames: []string{},
	}

	s := NewServer(planner)
	resultPB, err := s.ExtractVisFuncsInfo(context.Background(), req)
	require.Nil(t, err)
	assert.Equal(t, expected, resultPB)
}

func TestServerExtractVisFuncsInfo_FuncFilter(t *testing.T) {

	// Set up mocks.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	planner := mock_controller.NewMockPlanner(ctrl)

	visFuncsInfo := &scriptspb.VisFuncsInfo{}
	if err := proto.UnmarshalText(visFuncsInfoPBStr, visFuncsInfo); err != nil {
		t.Fatalf("Failed to unmarshal expected visfuncsinfo, err: %s", err)
	}
	expected := &scriptspb.VisFuncsInfo{}
	if err := proto.UnmarshalText(expectedVisFuncsInfoSingleFilter, expected); err != nil {
		t.Fatalf("Failed to unmarshal expected visfuncsinfo, err: %s", err)
	}

	planner.
		EXPECT().
		ExtractVisFuncsInfo(visFuncsQuery).
		Return(&scriptspb.VisFuncsInfoResult{
			Info:   visFuncsInfo,
			Status: &statuspb.Status{ErrCode: statuspb.OK},
		}, nil)

	req := &scriptmgrpb.ExtractVisFuncsInfoRequest{
		Script:    visFuncsQuery,
		FuncNames: []string{"f"},
	}

	s := NewServer(planner)
	resultPB, err := s.ExtractVisFuncsInfo(context.Background(), req)
	require.Nil(t, err)
	assert.Equal(t, expected, resultPB)
}

func TestServerExtractVisFuncsInfo_FuncFilter_Multiple(t *testing.T) {

	// Set up mocks.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	planner := mock_controller.NewMockPlanner(ctrl)

	visFuncsInfo := &scriptspb.VisFuncsInfo{}
	if err := proto.UnmarshalText(visFuncsInfoPBStr, visFuncsInfo); err != nil {
		t.Fatalf("Failed to unmarshal expected visfuncsinfo, err: %s", err)
	}
	expected := &scriptspb.VisFuncsInfo{}
	if err := proto.UnmarshalText(expectedVisFuncsInfoMultiFilter, expected); err != nil {
		t.Fatalf("Failed to unmarshal expected visfuncsinfo, err: %s", err)
	}

	planner.
		EXPECT().
		ExtractVisFuncsInfo(visFuncsQuery).
		Return(&scriptspb.VisFuncsInfoResult{
			Info:   visFuncsInfo,
			Status: &statuspb.Status{ErrCode: statuspb.OK},
		}, nil)

	req := &scriptmgrpb.ExtractVisFuncsInfoRequest{
		Script:    visFuncsQuery,
		FuncNames: []string{"f", "g"},
	}

	s := NewServer(planner)
	resultPB, err := s.ExtractVisFuncsInfo(context.Background(), req)
	require.Nil(t, err)
	assert.Equal(t, expected, resultPB)
}

func TestServerExtractVisFuncsInfo_FuncFilter_InvalidFuncName(t *testing.T) {
	// Set up mocks.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	planner := mock_controller.NewMockPlanner(ctrl)

	visFuncsInfo := &scriptspb.VisFuncsInfo{}
	if err := proto.UnmarshalText(visFuncsInfoPBStr, visFuncsInfo); err != nil {
		t.Fatalf("Failed to unmarshal expected visfuncsinfo, err: %s", err)
	}

	planner.
		EXPECT().
		ExtractVisFuncsInfo(visFuncsQuery).
		Return(&scriptspb.VisFuncsInfoResult{
			Info:   visFuncsInfo,
			Status: &statuspb.Status{ErrCode: statuspb.OK},
		}, nil)

	funcName := "this function doesn't exist"
	req := &scriptmgrpb.ExtractVisFuncsInfoRequest{
		Script:    visFuncsQuery,
		FuncNames: []string{funcName},
	}

	s := NewServer(planner)
	_, err := s.ExtractVisFuncsInfo(context.Background(), req)
	require.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("function '%s' was not found in script", funcName), err.Error())
}

func TestServerExtractVisFuncsInfo_CompilerError(t *testing.T) {
	// Set up mocks.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	planner := mock_controller.NewMockPlanner(ctrl)
	query := "this is an invalid query"
	fakeErrMsg := "compiler error"
	planner.
		EXPECT().
		ExtractVisFuncsInfo(query).
		Return(&scriptspb.VisFuncsInfoResult{
			Info: nil,
			Status: &statuspb.Status{
				ErrCode: statuspb.INVALID_ARGUMENT,
				Msg:     fakeErrMsg,
			},
		}, nil)

	req := &scriptmgrpb.ExtractVisFuncsInfoRequest{
		Script:    query,
		FuncNames: []string{},
	}

	s := NewServer(planner)
	_, err := s.ExtractVisFuncsInfo(context.Background(), req)
	require.NotNil(t, err)
	assert.Equal(t, fakeErrMsg, err.Error())
}
