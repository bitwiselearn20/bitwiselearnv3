"use client"
import React from 'react'
import Attempt from '@/component/attempt/Attempt'
import { useParams } from 'next/navigation'
const page = () => {
    const params = useParams();
  return (
    <div>
        <Attempt id={params.id as string} mode={"ASSESSMENT"} />
    </div>
  )
}

export default page