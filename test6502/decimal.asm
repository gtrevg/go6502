		.ORG $600

X 		.EQ 		$FE

START		
		LDX ABS:X,Y
		SED
		CLC
		LDA #$84
		ADC #$25
		BRK

TABLE		.BYTE		0x200
