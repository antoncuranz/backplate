#ifndef ARDUINO_INKPLATE10
#error "Wrong board selection for this example, please select Inkplate 10 in the boards menu."
#endif

#include "HTTPClient.h"
#include "Inkplate.h"
#include "WiFi.h"

#define uS_TO_MIN_FACTOR 60000000
#define TIME_TO_SLEEP    5

Inkplate display(INKPLATE_3BIT); // Create an object on Inkplate library and also set library into 1 Bit mode (BW)

const char ssid[] = "SNIP";
const char *password = "SNIP";

// use waveform5 from Inkplate10_Wavefrom_EEPROM_Programming (Arduino Library examples)

void setup()
{
  Serial.begin(115200);
  Serial.print("Connecting to WiFi...");

  // Connect to the WiFi network.
  WiFi.mode(WIFI_MODE_STA);
  WiFi.begin(ssid, password);
  while (WiFi.status() != WL_CONNECTED)
  {
    delay(500);
    Serial.print(".");
  }
  Serial.println("\nWiFi OK! Downloading...");

  display.begin(); // Init Inkplate library (you should call this function ONLY ONCE)
  display.setRotation(1); // Portrait
  display.clearDisplay(); // Clear frame buffer

  // Full color 24 bit images are large and take a long time to load, will take around 20 secs.
  HTTPClient http;
  // Set parameters to speed up the download process.
  http.getStream().setNoDelay(true);
  http.getStream().setTimeout(1);

  http.begin("http://TODO/consume");

  int httpCode = http.GET();
  if (httpCode == 200)
  {
    // Get the response length and make sure it is not 0.
    int32_t len = http.getSize();
    if (len > 0)
    {
      // REMEMBER! You can only use Windows Bitmap file with color depth of 1, 4, 8 or 24 bits with no compression!
      if (display.drawBitmapFromWeb(http.getStreamPtr(), 0, 0, len, false, false))
      {
        display.display();
      }
      else Serial.println("Image open error");
    }
    else Serial.println("Invalid response length");
  }
  else Serial.printf("HTTP error (status code %d)", httpCode);

  WiFi.mode(WIFI_OFF);

  Serial.println("Going to sleep now");
  delay(1000);
  Serial.flush();
  esp_sleep_enable_timer_wakeup(TIME_TO_SLEEP * uS_TO_MIN_FACTOR);
  esp_sleep_enable_ext0_wakeup(GPIO_NUM_36, LOW); // wake button
  esp_deep_sleep_start();
}

void loop()
{
  // Nothing...
}
