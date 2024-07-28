import { colors } from '@/constants/styles/colors'
import { TouchableOpacity, Text, StyleSheet } from 'react-native'
import { Ionicons } from '@expo/vector-icons'

interface ButtonProps {
  onPress: () => void
}

const MyShuttleButton = ({ onPress }: ButtonProps) => {
  return (
    <TouchableOpacity style={styles.button} onPress={onPress}>
      <Text style={styles.text}>My Shuttle</Text>
      <Ionicons name='chevron-forward' size={16} style={styles.icon} />
    </TouchableOpacity>
  )
}

const styles = StyleSheet.create({
  button: {
    padding: 16,
    borderRadius: 8,
    alignItems: 'center',
    justifyContent: 'center',
    display: 'flex',
    flexDirection: 'row',
    gap: 8,
    backgroundColor: colors.green,
  },
  icon: {
    color: colors.white,
  },
  text: {
    color: colors.white,
    fontFamily: 'Montserrat-Regular',
    textTransform: 'uppercase',
  },
})

export default MyShuttleButton
