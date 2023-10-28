from hypothesis import settings
from schemathesis import DataGenerationMethod
from . import schema, EXAMPLE_COUNT, UserAuth

@schema.parametrize(operation_id="post-user")
@settings(max_examples=EXAMPLE_COUNT)
def test_post_user(case):
    case.call_and_validate()

@schema.auth(UserAuth)
@schema.parametrize(operation_id="get-user", data_generation_methods=[DataGenerationMethod.positive]) # Not possible to do negative tests on this endpoint
@settings(max_examples=EXAMPLE_COUNT)
def test_get_user(case):
  case.call_and_validate()

@schema.auth(UserAuth)
@schema.parametrize(operation_id="patch-user")
@settings(max_examples=EXAMPLE_COUNT)
def test_patch_user(case):
		case.call_and_validate()

@schema.auth(UserAuth)
@schema.parametrize(operation_id="delete-user", data_generation_methods=[DataGenerationMethod.positive])
@settings(max_examples=EXAMPLE_COUNT)
def test_delete_user(case):
		case.call_and_validate()
