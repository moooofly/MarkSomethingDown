# GitHub Access Token

![octotree 被限速](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/Access%20Token%20problem%20when%20using%20octotree.png "octotree 被限速")

Octotree 基于 [GitHub API](https://developer.github.com/v3/) 获取 repository metadata 信息；默认情况下，它基于 **unauthenticated requests** 访问 GitHub API ；然而，存在两种 requests 必须被鉴权的情况：

- 访问一个 private repository 时
- 超过了[针对 unauthenticated requests 的速率限制](https://developer.github.com/v3/#rate-limiting)

当发生上述两种情况时，Octotree 将要求你提供 [GitHub personal access token](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/) ；如果你之前没有，可以先[创建一个](https://github.com/settings/tokens/new)，再将其复制粘贴到对应的文本框中；注意，最小权限访问至少要设置为 `public_repo` 和 `repo`（如果你需要能够访问 private repositories 的话）

![Cache/Pictures/new personal access token - 1](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/new%20personal%20access%20token%20-%201.png "Cache/Pictures/new personal access token - 1")

![Cache/Pictures/new personal access token - 2](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/new%20personal%20access%20token%20-%202.png "Cache/Pictures/new personal access token - 2")


- personal access token 的作用就和 OAuth access token 一样；
- 可被用于基于 HTTPS 的 git 模式中，从而避免使用密码；
- 可被用于基于 Basic Authentication 的 API 的鉴权；



----------


## [Rate Limiting](https://developer.github.com/v3/#rate-limiting)

任意一个 API 请求所返回的 HTTP headers 中都包含了当前的 rate limit 状态值：

```
curl -i https://api.github.com/users/octocat
HTTP/1.1 200 OK
Date: Mon, 01 Jul 2013 17:27:06 GMT
Status: 200 OK
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 56
X-RateLimit-Reset: 1372700873
```

| Header Name | Description |
| -- | -- | 
| **X-RateLimit-Limit** | 每小时允许发起的最大请求数量 |
| **X-RateLimit-Remaining** | 在当前 rate limit 窗口中仍允许发起的请求数量 |
| **X-RateLimit-Reset** | 重置当前 rate limit 时间窗口的时间周期（以 UTC epoch seconds 为单位） |

如果超出了上述 rate limit ，则会收到如下错误返回信息：

```
HTTP/1.1 403 Forbidden
Date: Tue, 20 Aug 2013 14:50:41 GMT
Status: 403 Forbidden
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1377013266

{
   "message": "API rate limit exceeded for xxx.xxx.xxx.xxx. (But here's the good news: Authenticated requests get a higher rate limit. Check out the documentation for more details.)",
   "documentation_url": "https://developer.github.com/v3/#rate-limiting"
}
```

## GitHub API rate limit

```shell
Error: GitHub API rate limit exceeded for 220.166.254.65. (But here's the good news: Authenticated requests get a higher rate limit. Check out the documentation for more details.)
Try again in 56 minutes 17 seconds, or create an personal access token:
  https://github.com/settings/tokens
and then set it as HOMEBREW_GITHUB_API_TOKEN.
```

解决办法：

- 注册、登录 https://github.com ；
- 访问 https://github.com/settings/tokens ；
- 在 `Personal settings -> Personal access tokens` 中点击 "Generate new token" 创建 token ；
- Make sure to copy your new personal access token now. You won't be able to see it again!
- Personal access tokens function like ordinary OAuth access tokens. They can be used instead of a password for Git over HTTPS, or can be used to authenticate to the API over Basic Authentication.
- 在 `~/.bashrc` 中添加 `HOMEBREW_GITHUB_API_TOKEN` 环境变量

```shell
if [ -f /usr/local/bin/brew ]; then
    export HOMEBREW_GITHUB_API_TOKEN=xxxxxxxxxx
fi
```


## [Creating a personal access token for the command line](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/)

> You can create a personal access token and use it in place of a password when performing Git operations over HTTPS with Git on the command line or the API.

在以下场景中，personal access token 被用于鉴权 GitHub 访问：

- 当你使用 [two-factor authentication](https://help.github.com/articles/about-two-factor-authentication) 时
- 基于 SAML single sign-on (SSO) 访问组织中受保护的内容时；Tokens used with organizations that use SAML SSO must be authorized.

### Creating a token

https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/#creating-a-token

### Using a token on the command line

一旦你拥有了 token ，你就可以在执行 Git operations over HTTPS 时使用 token 而不是密码；

例如，在命令行上你可以输入ß

```
git clone https://github.com/username/repo.git
Username: your_username
Password: your_token
```

Personal access tokens 只能被用于 HTTPS Git 操作；如果你的 repository 使用 SSH 形式的 remote URL ，那么你将需要[将 remote 从 SSH 切换成 HTTPS](https://help.github.com/articles/changing-a-remote-s-url/#switching-remote-urls-from-ssh-to-https) ；

如果你在使用过程中，没有被提示需要输入用户名和密码，则有可能你的密码信息已经被缓存在了你的电脑上；你可以通过[更新 Keychain](https://help.github.com/articles/updating-credentials-from-the-osx-keychain) ，将其中保存的旧密码替换成 token ；
