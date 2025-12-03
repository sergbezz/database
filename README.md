## Пример использования

```
> SET weather_2_pm cold_moscow_weather
OK
> GET weather_2_pm
cold_moscow_weather
> DEL weather_2_pm
OK
> GET weather_2_pm
Error: key not found
```
