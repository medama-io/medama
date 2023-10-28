from hypothesis import settings
from . import schema, EXAMPLE_COUNT, UserAuth

@schema.auth(UserAuth)
@schema.parametrize(operation_id="auth-login")
@settings(max_examples=EXAMPLE_COUNT)
def test_auth_login(case):
    case.call_and_validate()
