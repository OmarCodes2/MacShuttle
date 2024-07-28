import { colors } from '@/constants/styles/colors'
import { View, Text, StyleSheet } from 'react-native'

interface BottomSheetTitleProps {
  title: string
  subtitle: string
}

const BottomSheetTitle = ({ title, subtitle }: BottomSheetTitleProps) => {
  return (
    <View style={styles.container}>
      <Text style={styles.title}>{title}</Text>
      <Text style={styles.subtitle}>{subtitle}</Text>
    </View>
  )
}

export default BottomSheetTitle

const styles = StyleSheet.create({
  container: {
    display: 'flex',
    flexDirection: 'column',
    gap: 2,
  },
  title: {
    fontSize: 18,
    color: colors.white,
    fontFamily: 'Montserrat-Regular',
  },
  subtitle: {
    fontSize: 12,
    color: colors.white,
  },
})
