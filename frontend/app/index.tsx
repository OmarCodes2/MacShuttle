import React, { useState, useEffect, useRef } from 'react';
import { Text, View, StyleSheet, Button, ActivityIndicator } from 'react-native';
import * as Location from 'expo-location';
import axios from 'axios';

export default function Home() {
  const [location, setLocation] = useState<Location.LocationObject | null>(null);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const [tracking, setTracking] = useState(false);
  const [direction, setDirection] = useState('forward');
  const [startTime, setStartTime] = useState<number | null>(null);
  const intervalRef = useRef<NodeJS.Timeout | null>(null);

  useEffect(() => {
    (async () => {
      let { status } = await Location.requestForegroundPermissionsAsync();
      if (status !== 'granted') {
        setErrorMsg('Permission to access location was denied');
        return;
      }
    })();

    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
      }
    };
  }, []);

  const startTracking = () => {
    setTracking(true);
    const initialTime = Date.now();
    setStartTime(initialTime);
    intervalRef.current = setInterval(async () => {
      let currentLocation = await Location.getCurrentPositionAsync({});
      setLocation(currentLocation);

      const { latitude, longitude } = currentLocation.coords;
      const timestamp = Date.now() - initialTime; // Timestamp in milliseconds relative to start

      // Replace with your actual endpoint URL
      const endpointUrl = 'http://macshuttle-env.eba-thck352g.us-east-1.elasticbeanstalk.com/location';
      try {
        await axios.post(endpointUrl, {
          latitude,
          longitude,
          timestamp,
          direction,
        });
      } catch (error) {
        console.error('Error sending location data:', error);
      }
    }, 5000); // 5000ms (5 seconds)
  };

  const stopTracking = () => {
    setTracking(false);
    if (intervalRef.current) {
      clearInterval(intervalRef.current);
      intervalRef.current = null;
    }
  };

  const toggleDirection = () => {
    setDirection((prevDirection) => (prevDirection === 'forward' ? 'reverse' : 'forward'));
  };

  return (
    <View style={styles.container}>
      {errorMsg ? (
        <Text style={styles.error}>{errorMsg}</Text>
      ) : (
        location && (
          <Text style={styles.text}>
            Latitude: {location.coords.latitude}, Longitude: {location.coords.longitude}
          </Text>
        )
      )}
      <Button
        title={tracking ? 'Stop' : 'Start'}
        onPress={tracking ? stopTracking : startTracking}
      />
      {tracking && <ActivityIndicator size="large" color="#0000ff" />}
      <Button
        title={`Toggle Direction (Current: ${direction})`}
        onPress={toggleDirection}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: '#fff',
  },
  text: {
    fontSize: 24,
    color: '#000',
    margin: 10,
  },
  error: {
    fontSize: 18,
    color: 'red',
    margin: 10,
  },
});
