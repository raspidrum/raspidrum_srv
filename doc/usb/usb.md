

## Фильтрация списка USB устройств

**Необходимо исключать:**

 - hub: по bDeviceClass == 9 (Hub)
 - устройства ввода: В массиве `Configuration Descriptor` проверяем наличие `Interface Descriptor`, содержащего ТОЛЬКО `bInterfaceClass == 3 (Human Interface Device)`

**Поиск Audio устройств:**

  В массиве `Configuration Descriptor` проверяем наличие `Interface Descriptor`, содержащего `bInterfaceClass == 1 (Audio)`


**Наименование устройства:** iProduct

