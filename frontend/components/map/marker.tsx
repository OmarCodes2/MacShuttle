import { Marker } from 'react-native-maps'
import MarkerDefaultIcon from '../../assets/map/marker-default.svg'

interface MapMarkerProps {
  latitude: number
  longitude: number
}

const MapMarker = ({ latitude, longitude }: MapMarkerProps) => {
  return (
    <Marker coordinate={{ latitude, longitude }}>
      <MarkerDefaultIcon width={35} height={35} />
    </Marker>
  )
}

export default MapMarker
