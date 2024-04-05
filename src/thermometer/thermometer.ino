// NOTES: for getting list of connected devices: ls /dev/cu.*

// dependencies
#include <ESP8266WiFi.h>
#include <PubSubClient.h>
#include <Wire.h>
#include <GyverOLED.h>
#include "Adafruit_HTU21DF.h"
#include "TroykaRTC.h"
#include "secrets.h"

// consts
#define WAIT_TIMEOUT 10000
#define LOOP_DELAY 1000
#define DISPLAY_INIT_DELAY 1000
#define SERIAL_PORT 9600

// data
const char* title = SECRET_TITLE;
wl_status_t prev_status = wl_status_t::WL_IDLE_STATUS;

// wifi
const char* ssid = SECRET_SSID;
const char* pass = SECRET_PASS;

// mqtt
IPAddress mqtt_server(SECRET_SERVER_IP_OCTET_1, SECRET_SERVER_IP_OCTET_2, SECRET_SERVER_IP_OCTET_3, SECRET_SERVER_IP_OCTET_4);
const int mqtt_port = SECRET_MQTT_PORT;
const char* mqtt_user = SECRET_MQTT_USER;
const char* mqtt_pass = SECRET_MQTT_PASS;
const char* topic = SECRET_MQTT_TOPIC;
const char* mqtt_client_id = SECRET_MQTT_CLIENT_ID;

// services
WiFiClient wifi_client;
PubSubClient mqtt_client(mqtt_server, mqtt_port, wifi_client);
Adafruit_HTU21DF htu = Adafruit_HTU21DF();
GyverOLED< SSD1306_128x64, OLED_NO_BUFFER > oled;
RTC rtc_clock;

// init
void setup() {
  init_serial();
  init_clock();
  init_oled();
  init_sensor();
  try_init_wifi(WAIT_TIMEOUT);
}

// main loop
void loop() {
  float temperature = htu.readTemperature();
  float humidity = htu.readHumidity();
  uint32_t timestamp = get_timestamp();

  print_temperature(temperature);
  print_humidity(humidity);

  wl_status_t status = WiFi.status();
  print_wifi_status(status, prev_status);
  prev_status = status;

  delay(LOOP_DELAY);

  if (status != WL_CONNECTED) {
    try_init_wifi(WAIT_TIMEOUT);
  }

  if (WiFi.status() == WL_CONNECTED) {
    try_init_mqtt();
    send(temperature, humidity, timestamp);
  }
}

// functions
uint32_t get_timestamp() {
  rtc_clock.read();
  return rtc_clock.getUnixTime();
}

void send(float temperature, float humidity, uint32_t timestamp) {
  String message = "{ \"t\": #1, \"h\": #2, \"ts\": #3 }";
  message.replace("#1", String(temperature, 2));
  message.replace("#2", String(humidity, 2));
  message.replace("#3", String(timestamp));

  mqtt_client.publish(topic, message.c_str());
}

void print_temperature(float value) {
  oled.setCursor(0, 3);
  oled.print("Temp: ");
  oled.print(value);
}

void print_humidity(float value) {
  oled.setCursor(0, 4);
  oled.print("Humidity: ");
  oled.print(value);
}

void print_wifi_status(wl_status_t status, wl_status_t prev_status) {
  if (status == prev_status) {
    return;
  }

  oled.setCursor(0, 5);
  oled.print("                       ");
  oled.setCursor(0, 5);
  oled.print("WiFi: ");
  switch (status) {
    case 0:
      oled.print("Idle");
      break;
    case 1:
      oled.print("No SSID Avail");
      break;
    case 2:
      oled.print("Scan completed");
      break;
    case 3:
      oled.print("Connected");
      break;
    case 4:
      oled.print("Connect fail");
      break;
    case 5:
      oled.print("Connection lost");
      break;
    case 6:
      oled.print("Wrong pass");
      break;
    case 7:
      oled.print("Disconnected");
      break;
    case 255:
      oled.print("No shield");
      break;
    default:
      oled.print("Unknown");
      break;
  }
}

void try_init_wifi(int timeout) {
  const char* message = (String("Connecting to ") + ssid + String("...")).c_str();
  Serial.println(message);
  WiFi.begin(ssid, pass);

  if (WiFi.waitForConnectResult(timeout) != WL_CONNECTED) {
    return;
  }

  Serial.println("Connected to WiFi");
}

void try_init_mqtt() {
  if (!mqtt_client.connected()) {
    Serial.println("Connecting to MQTT server");
    if (mqtt_client.connect(mqtt_client_id, mqtt_user, mqtt_pass)) {
      Serial.println("Connected to MQTT server");
    } else {
      Serial.println("Could not connect to MQTT server");
    }
  } else {
    mqtt_client.loop();
  }
}

void init_clock() {
  rtc_clock.begin();
  rtc_clock.set(__TIMESTAMP__);
}

void init_serial() {
  Serial.begin(SERIAL_PORT);
}

void init_sensor() {
  if (!htu.begin()) {
    Serial.println("Couldn't find sensor!");
    while (1)
      ;
  }
}

void init_oled() {
  oled.init();
  oled.clear();
  oled.setScale(2);
  oled.home();
  oled.print(title);
  delay(DISPLAY_INIT_DELAY);
  oled.setScale(1);
}