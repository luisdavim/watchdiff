class Watchdiff < Formula
  desc "Watch a command and generate diffs from the output changes"
  homepage "https://github.com/luisdavim/watchdiff"
  url "https://github.com/luisdavim/watchdiff/releases/download/{{ .Version }}/watchdiff_darwin_x86_64.tar.gz"
  sha256 "{{ .Env.WATCHDIFF_DARWIN_X86_64_SHA256 }}"
  version "{{ .Version }}"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/luisdavim/watchdiff/releases/download/{{ .Version }}/watchdiff_darwin_x86_64.tar.gz"
    elsif Hardware::CPU.arm?
      url "https://github.com/luisdavim/watchdiff/releases/download/{{ .Version }}/watchdiff_darwin_arm64.tar.gz"
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/luisdavim/watchdiff/releases/download/{{ .Version }}/watchdiff_linux_x86_64.tar.gz"
    elsif Hardware::CPU.arm?
      url "https://github.com/luisdavim/watchdiff/releases/download/{{ .Version }}/watchdiff_linux_arm64.tar.gz"
    end
  end

  def install
    bin.install "watchdiff"
  end

  test do
    system "#{bin}/watchdiff", "-h"
  end
end
