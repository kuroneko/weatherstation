#include <DHT.h>
#include <DHT_U.h>

#define DHTPIN 2
#define DHTTYPE DHT11

#define DOORPIN 3

#define POLLDELAY 5000

DHT_Unified dht(DHTPIN, DHTTYPE);

void setup() {
  Serial.begin(9600);
  dht.begin();
  pinMode(DOORPIN, INPUT_PULLUP);
}

void loop() {
  static sensors_event_t  temp_event;
  static sensors_event_t  humidity_event;

  while (1) {
    dht.temperature().getEvent(&temp_event);
    dht.humidity().getEvent(&humidity_event);

    Serial.print("^");
    if(!isnan(temp_event.temperature)) {
      Serial.print(temp_event.temperature);
    } else {
      Serial.print("XXX");
    }
    Serial.print(":");
    if (!isnan(humidity_event.relative_humidity)) {
      Serial.print(humidity_event.relative_humidity);
    } else {
      Serial.print("XXX");
    }
    Serial.print(":");
    if (digitalRead(DOORPIN) == LOW) {
      Serial.print("0");
    } else {
      Serial.print("1");
    }
    Serial.println();

    delay(POLLDELAY);
  }
}
