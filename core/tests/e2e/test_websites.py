from hypothesis import settings
from schemathesis import DataGenerationMethod
from . import schema, EXAMPLE_COUNT, UserAuth


## POST /websites
@schema.parametrize(operation_id="post-websites")
@settings(max_examples=EXAMPLE_COUNT)
def test_post_websites(case):
    case.call_and_validate()


## GET /websites
@schema.auth(UserAuth)
@schema.parametrize(
    operation_id="^get-websites$",
    data_generation_methods=[DataGenerationMethod.positive],
)
@settings(max_examples=EXAMPLE_COUNT)
def test_get_websites(case):
    case.call_and_validate()


## GET /websites/{hostname}
@schema.auth(UserAuth)
@schema.parametrize(operation_id="^get-websites-id$")
@settings(max_examples=EXAMPLE_COUNT)
def test_get_websites_id(case):
    case.call_and_validate()


## PATCH /websites/{hostname}
@schema.auth(UserAuth)
@schema.parametrize(operation_id="patch-websites-id")
@settings(max_examples=EXAMPLE_COUNT)
def test_patch_websites_id(case):
    case.call_and_validate()


# DELETE /websites/{hostname}
@schema.auth(UserAuth)
@schema.parametrize(operation_id="delete-websites-id")
@settings(max_examples=EXAMPLE_COUNT)
def test_delete_websites_id(case):
    case.call_and_validate()
