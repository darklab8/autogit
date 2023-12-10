# Dev setup

## 📇 Dev Requirements

- golang
  - wget https://go.dev/dl/go1.20.3.linux-amd64.tar.gz
  - sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.20.3.linux-amd64.tar.gz
  - export PATH=$PATH:/usr/local/go/bin
- cobra generator https://github.com/spf13/cobra-cli/blob/main/README.md
- cobra guide https://github.com/spf13/cobra/blob/main/user_guide.md
- godoc
- add binary discovery for cobra-cli, godoc detection
  - `export PATH="$PATH:/usr/local/go/bin:$HOME/go/bin"`
- install latest stable autogit 😄
