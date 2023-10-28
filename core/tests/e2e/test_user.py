from hypothesis import settings
import requests
import schemathesis
from schemathesis import DataGenerationMethod

schema = schemathesis.from_uri("http://api:8080/openapi.yaml", data_generation_methods=[DataGenerationMethod.positive, DataGenerationMethod.negative], sanitize_output=False)
EXAMPLE_COUNT = 100

# Authentication
AUTH_ENDPOINT = "http://api:8080/auth/login"
CREATE_ENDPOINT = "http://api:8080/user"
TEST_EMAIL = "test@e2e.com"
TEST_PASSWORD = "test1234"

class UserAuth:
    def get(self, case, context):
      response = requests.post(
				CREATE_ENDPOINT,
				json={"email": TEST_EMAIL, "password": TEST_PASSWORD},
			)

			# User may already exist
      if response.status_code == 409:
        response = requests.post(
          AUTH_ENDPOINT,
					json={"email": TEST_EMAIL, "password": TEST_PASSWORD},
				)

        token = response.headers.get("Set-Cookie")
        return token

      if response.status_code != 201:
        raise Exception("Failed to create user" + response.text)

      token = response.headers.get("Set-Cookie")
      return token

    def set(self, case, data, context):
        case.headers = case.headers or {}
        case.headers["Cookie"] = data

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
