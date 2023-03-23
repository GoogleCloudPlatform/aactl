# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Aactl < Formula
  desc "Vulnerability management tool."
  homepage "https://github.com/GoogleCloudPlatform/aactl"
  version "0.4.7"
  license "Apache-2.0"

  on_macos do
    url "https://github.com/GoogleCloudPlatform/aactl/releases/download/v0.4.7/aactl_0.4.7_darwin_all"
    sha256 "11b938ace84b5532a3eb687943f4a87aa98eb4e3375332399a323afb654fdb89"

    def install
      bin.install "aactl_0.4.7_darwin_all" => "aactl"
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/GoogleCloudPlatform/aactl/releases/download/v0.4.7/aactl_0.4.7_linux_amd64"
      sha256 "d0fc82c0498c4db5793f70b3b29c6f6c792afb525b9f0c91da399ca02e5d2f1d"

      def install
        bin.install "aactl_0.4.7_linux_amd64" => "aactl"
      end
    end
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/GoogleCloudPlatform/aactl/releases/download/v0.4.7/aactl_0.4.7_linux_arm64"
      sha256 "630f954db4cdc8590196b6017a85f55c8586041ed5bc58cf7616484586056dd2"

      def install
        bin.install "aactl_0.4.7_linux_arm64" => "aactl"
      end
    end
  end

  test do
    system "#{bin}/aactl --version"
  end
end
