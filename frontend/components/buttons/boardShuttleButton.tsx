import { boardButtonType } from '@/constants/buttons/types'
import { colors } from '@/constants/styles/colors'
import { useMemo } from 'react'
import { TouchableOpacity, Text, StyleSheet } from 'react-native'

interface ButtonProps {
  status: boardButtonType
  onPress: () => void
}

const boardShuttleButton = ({ status, onPress }: ButtonProps) => {
  const buttonColor = useMemo(
    () => (status === 'exit' ? colors.green : colors.red),
    [status],
  )

  const buttonText = useMemo(
    () => (status === 'exit' ? 'Board' : 'Disembark'),
    [status],
  )

  return (
    <TouchableOpacity
      style={{ backgroundColor: buttonColor }}
      onPress={onPress}
    >
      <Text style={styles.text}>{buttonText}</Text>
    </TouchableOpacity>
  )
}

const styles = StyleSheet.create({
  button: {
    padding: 16,
    borderRadius: 8,
    alignItems: 'center',
  },
  text: {
    color: colors.white,
    fontFamily: 'Montserrat-Regular',
    textTransform: 'uppercase',
  },
})

export default boardShuttleButton
