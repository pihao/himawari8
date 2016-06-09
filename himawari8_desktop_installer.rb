#!/usr/bin/env ruby

APP_DIR = ENV['HOME'] + '/.himawari8_desktop'
PLIST_DIR = ENV['HOME'] + '/Library/LaunchAgents'
PLIST_FILE_NAME = 'com.pihao.himawari8_desktop.launchd.plist'
PROGRAM_FILE_NAME = 'himawari8_desktop.rb'
PLIST_FILE = PLIST_DIR + '/' + PLIST_FILE_NAME

def install
  puts "Generate plist file..."

  path = ENV['PATH']
  template = <<-XML
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>#{PLIST_FILE_NAME}</string>
    <key>EnvironmentVariables</key>
    <dict>
      <key>PATH</key>
      <string>#{path}</string>
    </dict>
    <key>ProgramArguments</key>
    <array>
        <string>#{APP_DIR}/#{PROGRAM_FILE_NAME}</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>StartInterval</key>
    <integer>600</integer>
    <key>StandardOutPath</key>
    <string>#{APP_DIR}/out.log</string>
    <key>StandardErrorPath</key>
    <string>#{APP_DIR}/err.log</string>
</dict>
</plist>
XML

  system "mkdir -p #{APP_DIR}"
  system "echo '#{template}' > #{PLIST_FILE}"

  puts "Config LaunchAgents..."

  system "cp #{File.dirname(__FILE__)}/#{PROGRAM_FILE_NAME} #{APP_DIR}"
  system "chmod u+x #{APP_DIR}/#{PROGRAM_FILE_NAME}"
  # system "launchctl setenv PATH $HOME/bin:/usr/local/bin:$PATH"
  system "launchctl load #{PLIST_FILE}"
  system "launchctl start #{PLIST_FILE}"

  puts "Install Complete."
end

def uninstall
  system "launchctl unload #{PLIST_FILE}"
  system("rm #{PLIST_FILE}")
  system("rm -rf #{APP_DIR}")
  puts "Uninstall Complete."
end


puts <<-TIP
Make sure you have installed `ImageMagick`:
    brew install imagemagick
    gem install mini_magick
TIP
puts "Please select: (i)nstall, (u)ninstall, (c)ancel:"
case gets.chomp
when 'i'
  install
when 'u'
  uninstall
else
  puts "Nothing happed."
end
