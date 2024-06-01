import React, { useState, useEffect, useRef } from 'react';
import { Text, View, StyleSheet, TouchableOpacity } from 'react-native';
import * as Location from 'expo-location';
import axios from 'axios';

export default function Home() {
  const [location, setLocation] = useState<Location.LocationObject | null>(null);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const [tracking, setTracking] = useState(false);
  const [direction, setDirection] = useState('forward');
  const [startTime, setStartTime] = useState<number | null>(null);
  const [runID, setRunID] = useState<number | null>(null);
  const intervalRef = useRef<NodeJS.Timeout | null>(null);
  const directionRef = useRef(direction);  // Added this line

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

  const startTracking = async () => {
    const retrieveID = `${process.env.EXPO_PUBLIC_API_URL}/startTracking`;
    try {
      const response = await axios.post(retrieveID);
      console.log('Retrieved RunID:', response.data.run_id);
      setRunID(response.data.run_id);
    } catch (error) {
      console.error('Error retrieving run ID:', error);
      return;
    }
    await Tracking();
  }

  const Tracking = () => {  // Changed from `const Tracking = async () => {` to `const Tracking = () => {`
    setTracking(true);
    const initialTime = Date.now();
    setStartTime(initialTime);
    intervalRef.current = setInterval(async () => {
      let currentLocation = await Location.getCurrentPositionAsync({});
      setLocation(currentLocation);

      const { latitude, longitude } = currentLocation.coords;
      const timestamp = Date.now() - initialTime;

      const endpointUrl = `${process.env.EXPO_PUBLIC_API_URL}/liveTracking`;
      console.log(endpointUrl)
      console.log("latitude is", latitude)

      try {
        console.log(directionRef.current)  // Changed from `direction` to `directionRef.current`
        await axios.post(endpointUrl, {
          latitude,
          longitude,
          timestamp,
          direction: directionRef.current,  // Changed from `direction` to `directionRef.current`
        });
      } catch (error) {
        console.error('Error sending location data:', error);
      }
    }, 500);
  };

  const stopTracking = () => {
    setTracking(false);
    if (intervalRef.current) {
      clearInterval(intervalRef.current);
      intervalRef.current = null;
    }
  };

  const toggleDirection = () => {
    setDirection((prevDirection) => {
      const newDirection = prevDirection === 'forward' ? 'reverse' : 'forward';
      directionRef.current = newDirection;  // Added this line
      return newDirection;
    });
  };

  useEffect(() => {
    console.log('Direction toggled to:', direction);
  }, [direction]);

  return (
    <View style={styles.container}>
      {errorMsg ? (
        <Text style={styles.error}>{errorMsg}</Text>
      ) : (
        <>
          {location && (
            <Text style={styles.text}>
              Latitude: {location.coords.latitude}, Longitude: {location.coords.longitude}
            </Text>
          )}
          {runID !== null && (
            <Text style={styles.text}>
              Current Run ID: {runID}
            </Text>
          )}
        </>
      )}
      <TouchableOpacity
        style={[styles.button, styles.startStopButton]}
        onPress={tracking ? stopTracking : startTracking}
      >
        <Text style={styles.buttonText}>{tracking ? 'Stop' : 'Start'}</Text>
      </TouchableOpacity>
      <TouchableOpacity style={styles.button} onPress={toggleDirection}>
        <Text style={styles.buttonText}>Toggle Direction (Current: {direction})</Text>
      </TouchableOpacity>
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
    margin: 20,
  },
  error: {
    fontSize: 18,
    color: 'red',
    margin: 10,
  },
  button: {
    backgroundColor: '#007bff',
    padding: 15,
    borderRadius: 25,
    margin: 10,
    width: '80%',
    alignItems: 'center',
  },
  startStopButton: {
    marginTop: 30,
  },
  buttonText: {
    color: '#fff',
    fontSize: 16,
  },
});
