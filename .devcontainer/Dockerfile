FROM golang:1.23

# Instala sudo, git, wget, unzip y crea el usuario arheanja con permisos de sudo
RUN apt-get update && apt-get install -y sudo git wget unzip \
    && useradd -ms /bin/bash arheanja --uid 1000 -G sudo \
    && echo "arheanja ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers \
    && wget https://releases.hashicorp.com/terraform/1.12.1/terraform_1.12.1_linux_amd64.zip \
    && unzip terraform_1.12.1_linux_amd64.zip \
    && mv terraform /usr/local/bin/ \
    && rm terraform_1.12.1_linux_amd64.zip

# Alias personalizados para todos los usuarios
RUN echo "\
alias update='sudo softwareupdate -i -a'\n\
alias c='clear'\n\
alias py='python3'\n\
alias v='vim'\n\
alias ll='ls -la'\n\
alias duh='du -sh'\n\
alias rmd='sudo find . -name \".DS_Store\" -depth -exec rm {} \;'\n\
alias gst='git status'\n\
alias gp='git pull'\n\
alias gpu='git push'\n\
alias gpo='git push origin'\n\
alias gcb='git checkout -b'\n\
alias gcm='git checkout main'\n\
alias gbd='git branch -D'\n\
alias gc='git commit -am'\n\
alias t='terraform'\n\
alias tu='terraform update'\n\
alias tf='terraform fmt'\n\
alias ti='terraform init'\n\
alias tp='terraform plan'\n\
alias ta='terraform apply'\n\
alias k='kubectl'\n\
alias g='get'\n\
" >> /etc/bash.bashrc

# Da permisos al usuario sobre el directorio de trabajo
RUN mkdir -p /go/src/app && chown -R arheanja:arheanja /go/src/app

USER arheanja

WORKDIR /go/src/app

CMD ["sleep", "infinity"]