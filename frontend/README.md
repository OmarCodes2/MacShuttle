# MacShuttle Frontend

This is an [Expo](https://expo.dev) project.

## Development

1. Install dependencies

   ```bash
   npm install
   ```

2. Start the app

   ```bash
    npx expo start
   ```

In the output, you'll find options to open the app in a

- [development build](https://docs.expo.dev/develop/development-builds/introduction/)
- [Android emulator](https://docs.expo.dev/workflow/android-studio-emulator/)
- [iOS simulator](https://docs.expo.dev/workflow/ios-simulator/)
- [Expo Go](https://expo.dev/go), a limited sandbox for trying out app development with Expo

## File structure

- `app` - main app
- `assets` - assets by components, divided by folder of concern
- `components` - components used by main app, divided by folder of concern
- `constants`- constant data used by components, divided by folder of concern

Folders of concern are consistent across assets, components, and constants. Note that this project uses [file-based routing](https://docs.expo.dev/router/introduction).

## Expo documentation

To see more about developing with Expo, look at the following resources:

- [Expo documentation](https://docs.expo.dev/): Learn fundamentals, or go into advanced topics with [guides](https://docs.expo.dev/guides)
- [Learn Expo tutorial](https://docs.expo.dev/tutorial/introduction/): Step-by-step tutorial to create a project that runs on Android, iOS, and the web
