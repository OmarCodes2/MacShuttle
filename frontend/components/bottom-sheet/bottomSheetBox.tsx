import React, { useRef } from 'react'
import { StyleSheet } from 'react-native'
import BottomSheet, { BottomSheetScrollView } from '@gorhom/bottom-sheet'
import { colors } from '@/constants/styles/colors'

const BottomSheetBox: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const bottomSheetRef = useRef<BottomSheet>(null)

  return (
    <BottomSheet
      ref={bottomSheetRef}
      index={0}
      snapPoints={['25%', '50%', '75%']}
      backgroundStyle={styles.background}
    >
      <BottomSheetScrollView contentContainerStyle={styles.contentContainer}>
        {children}
      </BottomSheetScrollView>
    </BottomSheet>
  )
}

export default BottomSheetBox

const styles = StyleSheet.create({
  background: {
    backgroundColor: colors.black,
  },
  contentContainer: {
    flex: 1,
    padding: 20,
  },
})
