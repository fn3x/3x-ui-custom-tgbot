#!/bin/bash

red='\033[0;31m'
green='\033[0;32m'
yellow='\033[0;33m'
plain='\033[0m'

cur_dir=$(pwd)
build_dir="build"

# check root
[[ $EUID -ne 0 ]] && echo -e "${red}Fatal error: ${plain} Please run this script with root privilege \n " && exit 1

# Check OS and set release variable
if [[ -f /etc/os-release ]]; then
    source /etc/os-release
    release=$ID
elif [[ -f /usr/lib/os-release ]]; then
    source /usr/lib/os-release
    release=$ID
else
    echo "Failed to check the system OS, please contact the author!" >&2
    exit 1
fi

if [ ! -f "$build_dir/x-ui" ]; then
  echo "Failed to install x-ui. Build the executable first! go build -o build/x-ui -v" >&2
  exit 1
fi

echo "The OS release is: $release"

arch() {
    case "$(uname -m)" in
    x86_64 | x64 | amd64) echo 'amd64' ;;
    i*86 | x86) echo '386' ;;
    armv8* | armv8 | arm64 | aarch64) echo 'arm64' ;;
    armv7* | armv7 | arm) echo 'armv7' ;;
    armv6* | armv6) echo 'armv6' ;;
    armv5* | armv5) echo 'armv5' ;;
    s390x) echo 's390x' ;;
    *) echo -e "${green}Unsupported CPU architecture! ${plain}" && rm -f install.sh && exit 1 ;;
    esac
}

echo "arch: $(arch)"

os_version=""
os_version=$(grep -i version_id /etc/os-release | cut -d \" -f2 | cut -d . -f1)

if [[ "${release}" == "arch" ]]; then
    echo "Your OS is Arch Linux"
elif [[ "${release}" == "parch" ]]; then
    echo "Your OS is Parch linux"
elif [[ "${release}" == "manjaro" ]]; then
    echo "Your OS is Manjaro"
elif [[ "${release}" == "armbian" ]]; then
    echo "Your OS is Armbian"
elif [[ "${release}" == "opensuse-tumbleweed" ]]; then
    echo "Your OS is OpenSUSE Tumbleweed"
elif [[ "${release}" == "centos" ]]; then
    if [[ ${os_version} -lt 8 ]]; then
        echo -e "${red} Please use CentOS 8 or higher ${plain}\n" && exit 1
    fi
elif [[ "${release}" == "ubuntu" || "${release}" == "pop" ]]; then
    if [[ ${os_version} -lt 20 ]]; then
        echo -e "${red} Please use Ubuntu 20 or higher version!${plain}\n" && exit 1
    fi
elif [[ "${release}" == "fedora" ]]; then
    if [[ ${os_version} -lt 36 ]]; then
        echo -e "${red} Please use Fedora 36 or higher version!${plain}\n" && exit 1
    fi
elif [[ "${release}" == "debian" ]]; then
    if [[ ${os_version} -lt 11 ]]; then
        echo -e "${red} Please use Debian 11 or higher ${plain}\n" && exit 1
    fi
elif [[ "${release}" == "almalinux" ]]; then
    if [[ ${os_version} -lt 9 ]]; then
        echo -e "${red} Please use AlmaLinux 9 or higher ${plain}\n" && exit 1
    fi
elif [[ "${release}" == "rocky" ]]; then
    if [[ ${os_version} -lt 9 ]]; then
        echo -e "${red} Please use Rocky Linux 9 or higher ${plain}\n" && exit 1
    fi
elif [[ "${release}" == "oracle" ]]; then
    if [[ ${os_version} -lt 8 ]]; then
        echo -e "${red} Please use Oracle Linux 8 or higher ${plain}\n" && exit 1
    fi
else
    echo -e "${red}Your operating system is not supported by this script.${plain}\n"
    echo "Please ensure you are using one of the following supported operating systems:"
    echo "- PopOs! 20.04+"
    echo "- Ubuntu 20.04+"
    echo "- Debian 11+"
    echo "- CentOS 8+"
    echo "- Fedora 36+"
    echo "- Arch Linux"
    echo "- Parch Linux"
    echo "- Manjaro"
    echo "- Armbian"
    echo "- AlmaLinux 9+"
    echo "- Rocky Linux 9+"
    echo "- Oracle Linux 8+"
    echo "- OpenSUSE Tumbleweed"
    exit 1

fi

gen_random_string() {
    local length="$1"
    local random_string=$(LC_ALL=C tr -dc 'a-zA-Z0-9' </dev/urandom | fold -w "$length" | head -n 1)
    echo "$random_string"
}

# This function will be called when user installed x-ui out of security
config_after_install() {
    echo -e "${yellow}Install/update finished! For security it's recommended to modify panel settings ${plain}"
    read -p "Would you like to customize the panel settings? (If not, random settings will be applied) [y/n]: " config_confirm
    if [[ "${config_confirm}" == "y" || "${config_confirm}" == "Y" ]]; then
        read -p "Please set up your username: " config_account
        echo -e "${yellow}Your username will be: ${config_account}${plain}"
        read -p "Please set up your password: " config_password
        echo -e "${yellow}Your password will be: ${config_password}${plain}"
        read -p "Please set up the panel port: " config_port
        echo -e "${yellow}Your panel port is: ${config_port}${plain}"
        read -p "Please set up yookassa shop id: " shop_id
        echo -e "${yellow}Your shop id is: ${shop_id}${plain}"
        read -p "Please set up yookassa API key: " api_key
        echo -e "${yellow}Your API key is: ${api_key}${plain}"
        read -p "Please set up email for receipts: " email
        echo -e "${yellow}Your email for receipts is: ${email}${plain}"
        read -p "Please set up the port for youkassa webhooks: " webhook_port 
        echo -e "${yellow}Port for yookassa webhooks: ${webhook_port}${plain}"
        read -p "Please set up the web base path (ip:port/webbasepath/): " config_webBasePath
        echo -e "${yellow}Your web base path is: ${config_webBasePath}${plain}"
        echo -e "${yellow}Initializing, please wait...${plain}"
        /usr/local/x-ui/x-ui setting -username ${config_account} -password ${config_password}
        echo -e "${yellow}Account name and password set successfully!${plain}"
        /usr/local/x-ui/x-ui setting -port ${config_port}
        echo -e "${yellow}Panel port set successfully!${plain}"
        /usr/local/x-ui/x-ui setting -webBasePath ${config_webBasePath}
        echo -e "${yellow}Web base path set successfully!${plain}"
        /usr/local/x-ui/x-ui setting -shopId ${shop_id}
        /usr/local/x-ui/x-ui setting -apiKey ${api_key}
        echo -e "${yellow}Yookassa auth set successfully!${plain}"
        /usr/local/x-ui/x-ui setting -webhookPort ${webhook_port}
        echo -e "${yellow}Yookassa webhook port set successfully!${plain}"
        /usr/local/x-ui/x-ui setting -email ${email}
        echo -e "${yellow}Email for receipts set successfully!${plain}"
    else
        echo -e "${red}Cancel...${plain}"
        if [[ ! -f "/etc/x-ui/x-ui.db" ]]; then
            local usernameTemp=$(head -c 6 /dev/urandom | base64)
            local passwordTemp=$(head -c 6 /dev/urandom | base64)
            local webBasePathTemp=$(gen_random_string 10)
            /usr/local/x-ui/x-ui setting -username ${usernameTemp} -password ${passwordTemp} -webBasePath ${webBasePathTemp}
            echo -e "This is a fresh installation, will generate random login info for security concerns:"
            echo -e "###############################################"
            echo -e "${green}Username: ${usernameTemp}${plain}"
            echo -e "${green}Password: ${passwordTemp}${plain}"
            echo -e "${green}WebBasePath: ${webBasePathTemp}${plain}"
            echo -e "###############################################"
            echo -e "${yellow}If you forgot your login info, you can type "x-ui settings" to check after installation${plain}"
        else
            echo -e "${yellow}This is your upgrade, will keep old settings. If you forgot your login info, you can type "x-ui settings" to check${plain}"
        fi
    fi
    /usr/local/x-ui/x-ui migrate
}

install_base() {
    case "${release}" in
    ubuntu | debian | armbian | pop)
        apt-get update && apt-get install -y -q wget curl tar tzdata
        ;;
    centos | almalinux | rocky | oracle)
        yum -y update && yum install -y -q wget curl tar tzdata
        ;;
    fedora)
        dnf -y update && dnf install -y -q wget curl tar tzdata
        ;;
    arch | manjaro | parch)
        pacman -Syu && pacman -Syu --noconfirm wget curl tar tzdata
        ;;
    opensuse-tumbleweed)
        zypper refresh && zypper -q install -y wget curl tar timezone
        ;;
    *)
        apt-get update && apt install -y -q wget curl tar tzdata
        ;;
    esac
}

install_x-ui() {
    echo -e "Beginning to install x-ui $1"
    if [[ -e /usr/local/x-ui/ ]]; then
        systemctl stop x-ui
        rm /usr/local/x-ui/ -rf
    fi

    local CGO_ENABLED=1
    local GOOS=linux
    local architecture=arch
    local GOARCH=arch
    if [ $(arch) == "arm64" ]; then
      local GOARCH=arm64
      local CC=aarch64-linux-gnu-gcc
    elif [ $(arch) == "armv7" ]; then
      local GOARCH=arm
      local GOARM=7
      local CC=arm-linux-gnueabihf-gcc
    elif [ $(arch) == "armv6" ]; then
      local GOARCH=arm
      local GOARM=6
      local CC=arm-linux-gnueabihf-gcc
    elif [ $(arch) == "386" ]; then
      local GOARCH=386
      local CC=i686-linux-gnu-gcc
    elif [ $(arch) == "armv5" ]; then
      local GOARCH=arm
      local GOARM=5
      local CC=arm-linux-gnueabi-gcc
    elif [ $(arch) == "s390x" ]; then
      local GOARCH=s390x
      local CC=s390x-linux-gnu-gcc
    fi
    
    mkdir /usr/local/x-ui

    cp $build_dir/x-ui /usr/bin/x-ui

    mkdir -p $build_dir/bin
    cd $build_dir/bin
    
    # Download dependencies
    Xray_URL="https://github.com/XTLS/Xray-core/releases/download/v1.8.23/"
    if [ $(arch) == "amd64" ]; then
      wget ${Xray_URL}Xray-linux-64.zip
      unzip Xray-linux-64.zip
      rm -f Xray-linux-64.zip
    elif [ $(arch) == "arm64" ]; then
      wget ${Xray_URL}Xray-linux-arm64-v8a.zip
      unzip Xray-linux-arm64-v8a.zip
      rm -f Xray-linux-arm64-v8a.zip
    elif [ $(arch) == "armv7" ]; then
      wget ${Xray_URL}Xray-linux-arm32-v7a.zip
      unzip Xray-linux-arm32-v7a.zip
      rm -f Xray-linux-arm32-v7a.zip
    elif [ $(arch) == "armv6" ]; then
      wget ${Xray_URL}Xray-linux-arm32-v6.zip
      unzip Xray-linux-arm32-v6.zip
      rm -f Xray-linux-arm32-v6.zip
    elif [ $(arch) == "386" ]; then
      wget ${Xray_URL}Xray-linux-32.zip
      unzip Xray-linux-32.zip
      rm -f Xray-linux-32.zip
    elif [ $(arch) == "armv5" ]; then
      wget ${Xray_URL}Xray-linux-arm32-v5.zip
      unzip Xray-linux-arm32-v5.zip
      rm -f Xray-linux-arm32-v5.zip
    elif [ $(arch) == "s390x" ]; then
      wget ${Xray_URL}Xray-linux-s390x.zip
      unzip Xray-linux-s390x.zip
      rm -f Xray-linux-s390x.zip
    fi
    rm -f geoip.dat geosite.dat
    wget https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geoip.dat
    wget https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat
    wget -O geoip_IR.dat https://github.com/chocolate4u/Iran-v2ray-rules/releases/latest/download/geoip.dat
    wget -O geosite_IR.dat https://github.com/chocolate4u/Iran-v2ray-rules/releases/latest/download/geosite.dat
    wget -O geoip_VN.dat https://github.com/vuong2023/vn-v2ray-rules/releases/latest/download/geoip.dat
    wget -O geosite_VN.dat https://github.com/vuong2023/vn-v2ray-rules/releases/latest/download/geosite.dat
    mv xray xray-linux-$(arch)

    cd ..

    # Check the system's architecture and rename the file accordingly
    if [[ $(arch) == "armv5" || $(arch) == "armv6" || $(arch) == "armv7" ]]; then
        mv bin/xray-linux-$(arch) bin/xray-linux-arm
        chmod +x bin/xray-linux-arm
    fi

    cd ..

    cp -r $build_dir/* /usr/local/x-ui
    cp -f local.x-ui.service /etc/systemd/system/x-ui.service
    cp local_x-ui.sh /usr/local/x-ui/x-ui.sh

    chmod +x /usr/local/x-ui/x-ui.sh
    chmod +x /usr/bin/x-ui
    config_after_install

    systemctl daemon-reload
    systemctl enable x-ui
    systemctl start x-ui
    echo -e "${green}x-ui local_version${plain} installation finished, it is running now..."
    echo -e ""
    echo -e "x-ui control menu usages: "
    echo -e "----------------------------------------------"
    echo -e "SUBCOMMANDS:"
    echo -e "x-ui              - Admin Management Script"
    echo -e "x-ui start        - Start"
    echo -e "x-ui stop         - Stop"
    echo -e "x-ui restart      - Restart"
    echo -e "x-ui status       - Current Status"
    echo -e "x-ui settings     - Current Settings"
    echo -e "x-ui enable       - Enable Autostart on OS Startup"
    echo -e "x-ui disable      - Disable Autostart on OS Startup"
    echo -e "x-ui log          - Check logs"
    echo -e "x-ui banlog       - Check Fail2ban ban logs"
    echo -e "x-ui update       - Update"
    echo -e "x-ui custom       - custom version"
    echo -e "x-ui install      - Install"
    echo -e "x-ui uninstall    - Uninstall"
    echo -e "----------------------------------------------"
}

echo -e "${green}Running...${plain}"
install_base
install_x-ui $1
