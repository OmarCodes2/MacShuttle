import { Marker } from 'react-native-maps'
import MarkerDefaultIcon from '../../assets/map/marker-default.svg'

interface MapMarkerProps {
  latitude: number
  longitude: number
}

const MapMarker = ({ latitude, longitude }: MapMarkerProps) => {
  return (
    <Marker
      coordinate={{ latitude, longitude }}
      title='My Marker'
      description='This is a description of the marker'
    >
      <MarkerDefaultIcon width={30} height={30} />
    </Marker>
  )
}

export default MapMarker
