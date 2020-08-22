package main

        // A | addr
        func ORA( addr uint16)
        {
            cycleLock.enterCycle();
            A |= getByteFromMemory(address, lockToCycle: false);

            setProcessorStatusBit(ProcessorStatus.Z, isSet: (A == 0));
            setProcessorStatusBit(ProcessorStatus.N, isSet: ((A & (byte)ProcessorStatus.N) != 0));
            cycleLock.exitCycle();
        }