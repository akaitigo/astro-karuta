"use client";

import { useCallback, useEffect, useRef, useState } from "react";

interface UseAudioReturn {
  speak: (text: string) => void;
  stop: () => void;
  isSpeaking: boolean;
  isSupported: boolean;
}

/**
 * useAudio provides text-to-speech using the Web SpeechSynthesis API.
 * Prefers Japanese voice when available.
 */
export function useAudio(): UseAudioReturn {
  const [isSpeaking, setIsSpeaking] = useState(false);
  const [isSupported, setIsSupported] = useState(false);
  const japaneseVoiceRef = useRef<SpeechSynthesisVoice | null>(null);
  const synthRef = useRef<SpeechSynthesis | null>(null);

  useEffect(() => {
    if (typeof window === "undefined" || !window.speechSynthesis) {
      setIsSupported(false);
      return;
    }

    const synth = window.speechSynthesis;
    synthRef.current = synth;
    setIsSupported(true);

    const loadVoices = () => {
      const voices = synth.getVoices();
      const ja = voices.find(
        (v) => v.lang.startsWith("ja") && v.localService,
      );
      japaneseVoiceRef.current =
        ja ?? voices.find((v) => v.lang.startsWith("ja")) ?? null;
    };

    loadVoices();

    // Voices may load asynchronously in some browsers
    synth.addEventListener("voiceschanged", loadVoices);
    return () => {
      synth.removeEventListener("voiceschanged", loadVoices);
    };
  }, []);

  const stop = useCallback(() => {
    synthRef.current?.cancel();
    setIsSpeaking(false);
  }, []);

  const speak = useCallback(
    (text: string) => {
      if (!synthRef.current) return;

      // Cancel any ongoing speech
      stop();

      const utterance = new SpeechSynthesisUtterance(text);
      utterance.lang = "ja-JP";
      utterance.rate = 0.9;
      utterance.pitch = 1.0;

      if (japaneseVoiceRef.current) {
        utterance.voice = japaneseVoiceRef.current;
      }

      utterance.onstart = () => setIsSpeaking(true);
      utterance.onend = () => setIsSpeaking(false);
      utterance.onerror = () => setIsSpeaking(false);

      synthRef.current.speak(utterance);
    },
    [stop],
  );

  return { speak, stop, isSpeaking, isSupported };
}
