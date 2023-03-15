package elasticache_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfelasticache "github.com/hashicorp/terraform-provider-aws/internal/service/elasticache"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccElastiCacheUser_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var user elasticache.User
	rName := sdkacctest.RandomWithPrefix("tf-acc")
	resourceName := "aws_elasticache_user.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, elasticache.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUserDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists(ctx, resourceName, &user),
					resource.TestCheckResourceAttr(resourceName, "user_id", rName),
					resource.TestCheckResourceAttr(resourceName, "no_password_required", "false"),
					resource.TestCheckResourceAttr(resourceName, "user_name", "username1"),
					resource.TestCheckResourceAttr(resourceName, "engine", "redis"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"no_password_required",
					"passwords",
				},
			},
		},
	})
}

func TestAccElastiCacheUser_password_auth_mode(t *testing.T) {
	ctx := acctest.Context(t)
	var user elasticache.User
	rName := sdkacctest.RandomWithPrefix("tf-acc")
	resourceName := "aws_elasticache_user.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, elasticache.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUserDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfigWithPasswordAuthMode_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists(ctx, resourceName, &user),
					resource.TestCheckResourceAttr(resourceName, "user_id", rName),
					resource.TestCheckResourceAttr(resourceName, "user_name", "username1"),
					resource.TestCheckResourceAttr(resourceName, "engine", "redis"),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode.0.password_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode.0.passwords.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "authentication_mode.0.passwords.*", "aaaaaaaaaaaaaaaa"),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode.0.type", "password"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"authentication_mode.0.passwords.#",
					"authentication_mode.0.passwords.0",
					"no_password_required",
				},
			},
		},
	})
}

func TestAccElastiCacheUser_iam_auth_mode(t *testing.T) {
	ctx := acctest.Context(t)
	var user elasticache.User
	rName := sdkacctest.RandomWithPrefix("tf-acc")
	resourceName := "aws_elasticache_user.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, elasticache.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUserDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfigWithIamAuthMode_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists(ctx, resourceName, &user),
					resource.TestCheckResourceAttr(resourceName, "user_id", rName),
					resource.TestCheckResourceAttr(resourceName, "user_name", rName),
					resource.TestCheckResourceAttr(resourceName, "engine", "redis"),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode.0.type", "iam"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"no_password_required",
				},
			},
		},
	})
}

func TestAccElastiCacheUser_update(t *testing.T) {
	ctx := acctest.Context(t)
	var user elasticache.User
	rName := sdkacctest.RandomWithPrefix("tf-acc")
	resourceName := "aws_elasticache_user.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, elasticache.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUserDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists(ctx, resourceName, &user),
				),
			},
			{
				Config: testAccUserConfig_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists(ctx, resourceName, &user),
					resource.TestCheckResourceAttr(resourceName, "access_string", "on ~* +@all"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"no_password_required",
					"passwords",
				},
			},
		},
	})
}

func TestAccElastiCacheUser_update_password_auth_mode(t *testing.T) {
	ctx := acctest.Context(t)
	var user elasticache.User
	rName := sdkacctest.RandomWithPrefix("tf-acc")
	resourceName := "aws_elasticache_user.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, elasticache.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUserDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfigWithPasswordAuthMode_twoPasswords(rName, "aaaaaaaaaaaaaaaa", "bbbbbbbbbbbbbbbb"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists(ctx, resourceName, &user),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode.0.password_count", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"authentication_mode.0.passwords",
					"no_password_required",
				},
			},
			{
				Config: testAccUserConfigWithPasswordAuthMode_onePassword(rName, "aaaaaaaaaaaaaaaa"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists(ctx, resourceName, &user),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode.0.password_count", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"authentication_mode.0.passwords",
					"no_password_required",
				},
			},
			{
				Config: testAccUserConfigWithPasswordAuthMode_twoPasswords(rName, "cccccccccccccccc", "dddddddddddddddd"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists(ctx, resourceName, &user),
					resource.TestCheckResourceAttr(resourceName, "authentication_mode.0.password_count", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"authentication_mode.0.passwords",
					"no_password_required",
				},
			},
		},
	})
}

func TestAccElastiCacheUser_tags(t *testing.T) {
	ctx := acctest.Context(t)
	var user elasticache.User
	rName := sdkacctest.RandomWithPrefix("tf-acc")
	resourceName := "aws_elasticache_user.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, elasticache.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUserDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig_tags(rName, "tagKey", "tagVal"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists(ctx, resourceName, &user),
					resource.TestCheckResourceAttr(resourceName, "user_id", rName),
					resource.TestCheckResourceAttr(resourceName, "no_password_required", "false"),
					resource.TestCheckResourceAttr(resourceName, "user_name", "username1"),
					resource.TestCheckResourceAttr(resourceName, "engine", "redis"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.tagKey", "tagVal"),
				),
			},
			{
				Config: testAccUserConfig_tags(rName, "tagKey", "tagVal2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists(ctx, resourceName, &user),
					resource.TestCheckResourceAttr(resourceName, "user_id", rName),
					resource.TestCheckResourceAttr(resourceName, "no_password_required", "false"),
					resource.TestCheckResourceAttr(resourceName, "user_name", "username1"),
					resource.TestCheckResourceAttr(resourceName, "engine", "redis"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.tagKey", "tagVal2"),
				),
			},
			{
				Config: testAccUserConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists(ctx, resourceName, &user),
					resource.TestCheckResourceAttr(resourceName, "user_id", rName),
					resource.TestCheckResourceAttr(resourceName, "no_password_required", "false"),
					resource.TestCheckResourceAttr(resourceName, "user_name", "username1"),
					resource.TestCheckResourceAttr(resourceName, "engine", "redis"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccElastiCacheUser_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var user elasticache.User
	rName := sdkacctest.RandomWithPrefix("tf-acc")
	resourceName := "aws_elasticache_user.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, elasticache.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUserDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists(ctx, resourceName, &user),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfelasticache.ResourceUser(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckUserDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error { return testAccCheckUserDestroyWithProvider(ctx)(s, acctest.Provider) }
}

func testAccCheckUserDestroyWithProvider(ctx context.Context) acctest.TestCheckWithProviderFunc {
	return func(s *terraform.State, provider *schema.Provider) error {
		conn := provider.Meta().(*conns.AWSClient).ElastiCacheConn()

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_elasticache_user" {
				continue
			}

			user, err := tfelasticache.FindUserByID(ctx, conn, rs.Primary.ID)

			if tfawserr.ErrCodeEquals(err, elasticache.ErrCodeUserNotFoundFault) || tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			if user != nil {
				return fmt.Errorf("ElastiCache User (%s) still exists", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccCheckUserExists(ctx context.Context, n string, v *elasticache.User) resource.TestCheckFunc {
	return testAccCheckUserExistsWithProvider(ctx, n, v, func() *schema.Provider { return acctest.Provider })
}

func testAccCheckUserExistsWithProvider(ctx context.Context, n string, v *elasticache.User, providerF func() *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ElastiCache User ID is set")
		}

		provider := providerF()
		conn := provider.Meta().(*conns.AWSClient).ElastiCacheConn()
		resp, err := tfelasticache.FindUserByID(ctx, conn, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("ElastiCache User (%s) not found: %w", rs.Primary.ID, err)
		}

		*v = *resp

		return nil
	}
}

func testAccUserConfig_basic(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptIn(), fmt.Sprintf(`
resource "aws_elasticache_user" "test" {
  user_id       = %[1]q
  user_name     = "username1"
  access_string = "on ~app::* -@all +@read +@hash +@bitmap +@geo -setbit -bitfield -hset -hsetnx -hmset -hincrby -hincrbyfloat -hdel -bitop -geoadd -georadius -georadiusbymember"
  engine        = "REDIS"
  passwords     = ["password123456789"]
}
`, rName))
}

func testAccUserConfigWithPasswordAuthMode_basic(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptIn(), fmt.Sprintf(`
resource "aws_elasticache_user" "test" {
  user_id       = %[1]q
  user_name     = "username1"
  access_string = "on ~app::* -@all +@read +@hash +@bitmap +@geo -setbit -bitfield -hset -hsetnx -hmset -hincrby -hincrbyfloat -hdel -bitop -geoadd -georadius -georadiusbymember"
  engine        = "REDIS"

  authentication_mode {
    type      = "password"
    passwords = ["aaaaaaaaaaaaaaaa"]
  }
}
`, rName))
}

func testAccUserConfigWithIamAuthMode_basic(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptIn(), fmt.Sprintf(`
resource "aws_elasticache_user" "test" {
  user_id       = %[1]q
  user_name     = %[1]q
  access_string = "on ~app::* -@all +@read +@hash +@bitmap +@geo -setbit -bitfield -hset -hsetnx -hmset -hincrby -hincrbyfloat -hdel -bitop -geoadd -georadius -georadiusbymember"
  engine        = "REDIS"

  authentication_mode {
    type      = "iam"
  }
}
`, rName))
}

func testAccUserConfigWithNoPassRequiredAuthMode_basic(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptIn(), fmt.Sprintf(`
resource "aws_elasticache_user" "test" {
  user_id       = %[1]q
  user_name     = "username1"
  access_string = "on ~app::* -@all +@read +@hash +@bitmap +@geo -setbit -bitfield -hset -hsetnx -hmset -hincrby -hincrbyfloat -hdel -bitop -geoadd -georadius -georadiusbymember"
  engine        = "REDIS"

  authentication_mode {
    type      = "no-password-required"
  }
}
`, rName))
}

func testAccUserConfig_update(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptIn(), fmt.Sprintf(`
resource "aws_elasticache_user" "test" {
  user_id       = %[1]q
  user_name     = "username1"
  access_string = "on ~* +@all"
  engine        = "REDIS"
  passwords     = ["password234567891", "password345678912"]
}
`, rName))
}

func testAccUserConfigWithPasswordAuthMode_twoPasswords(rName string, password1 string, password2 string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptIn(), fmt.Sprintf(`
resource "aws_elasticache_user" "test" {
  user_id       = %[1]q
  user_name     = "username1"
  access_string = "on ~app::* -@all +@read +@hash +@bitmap +@geo -setbit -bitfield -hset -hsetnx -hmset -hincrby -hincrbyfloat -hdel -bitop -geoadd -georadius -georadiusbymember"
  engine        = "REDIS"

  authentication_mode {
    type      = "password"
    passwords = [%[2]q, %[3]q]
  }	
}
`, rName, password1, password2))
}

func testAccUserConfigWithPasswordAuthMode_onePassword(rName string, password string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptIn(), fmt.Sprintf(`
resource "aws_elasticache_user" "test" {
  user_id       = %[1]q
  user_name     = "username1"
  access_string = "on ~app::* -@all +@read +@hash +@bitmap +@geo -setbit -bitfield -hset -hsetnx -hmset -hincrby -hincrbyfloat -hdel -bitop -geoadd -georadius -georadiusbymember"
  engine        = "REDIS"

  authentication_mode {
    type      = "password"
    passwords = [%[2]q]
  }	
}
`, rName, password))
}

func testAccUserConfig_tags(rName, tagKey, tagValue string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptIn(), fmt.Sprintf(`
resource "aws_elasticache_user" "test" {
  user_id       = %[1]q
  user_name     = "username1"
  access_string = "on ~app::* -@all +@read +@hash +@bitmap +@geo -setbit -bitfield -hset -hsetnx -hmset -hincrby -hincrbyfloat -hdel -bitop -geoadd -georadius -georadiusbymember"
  engine        = "REDIS"
  passwords     = ["password123456789"]

  tags = {
    %[2]s = %[3]q
  }
}
`, rName, tagKey, tagValue))
}
