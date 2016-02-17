#!/usr/bin/env ruby

APP_DIR = ENV['HOME'] + '/.himawari8_desktop'


require 'json'
require 'net/http'
require "mini_magick"
require 'open-uri'
require 'uri'

def generateImage
  config_url = 'http://himawari8.nict.go.jp/img/D531106/latest.json'
  base_img_url = 'http://himawari8-dl.nict.go.jp/himawari8/img/D531106/2d/550/'
  # http://himawari8-dl.nict.go.jp/himawari8/img/D531106/2d/550/2016/01/08/035000_0_0.png

  #latest_date = JSON.parse(Net::HTTP.get(URI(config_url)))['date']
  latest_date = JSON.parse(open(config_url){|f|f.read})['date']
  file_name = APP_DIR + '/himawari8_' + latest_date.split(/[-| |:]/)[4] + '.png'
  format_date = latest_date.gsub("-", "/").gsub(" ", "/").gsub(":", "")

  urls = [
    base_img_url + format_date + '_0_0.png',
    base_img_url + format_date + '_1_0.png',
    base_img_url + format_date + '_0_1.png',
    base_img_url + format_date + '_1_1.png'
  ]

  MiniMagick::Tool::Convert.new do |convert|
    convert << "-size" << "1100x1100" << "xc:none"
    convert << "-strip"
    convert << urls[0] << "-geometry" << "+0+0" << "-composite"
    convert << urls[1] << "-geometry" << "+550+0" << "-composite"
    convert << urls[2] << "-geometry" << "+0+550" << "-composite"
    convert << urls[3] << "-geometry" << "+550+550" << "-composite"
    convert << file_name
  end

  file_name
end

def osascript(script)
  system 'osascript', *script.split(/\n/).map { |line| ['-e', line] }.flatten
end

def setDesktop(image)
  scpt = <<-SCPT
tell application "System Events"
    set desktopCount to count of desktops
    repeat with i from 1 to desktopCount
        tell desktop i
            set picture to "#{image}"
        end tell
    end repeat
end tell
SCPT
  osascript(scpt)
end

def debug
  system "echo '- - - - - - - - - - - -'"
  system "date"
  system "echo $USER"
  system "echo $PATH"
  system "ruby -v"
  system "gem list | ag mini_magick"
end

# debug
setDesktop generateImage
