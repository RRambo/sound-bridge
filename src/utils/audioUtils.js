// Audio utility for playing alert sounds

// Play alert sound
export const playAlertSound = () => {
  try {
    // Method 1: Try to play custom audio file
    const audio = new Audio('./public/alert.wav');
    audio.volume = 1.0; // 100% volume
    audio.play().catch((error) => {
      console.warn('Failed to play custom alert sound:', error);
      // Fallback: Use browser beep (Web Audio API)
      playBeep();
    });
  } catch (error) {
    console.error('Error playing alert sound:', error);
    // Fallback to beep
    playBeep();
  }
};

// Fallback: Generate beep sound using Web Audio API
const playBeep = () => {
  try {
    const audioContext = new (window.AudioContext || window.webkitAudioContext)();
    const oscillator = audioContext.createOscillator();
    const gainNode = audioContext.createGain();

    oscillator.connect(gainNode);
    gainNode.connect(audioContext.destination);

    oscillator.frequency.value = 800; // Frequency in Hz
    oscillator.type = 'sine';

    gainNode.gain.setValueAtTime(0.3, audioContext.currentTime);
    gainNode.gain.exponentialRampToValueAtTime(0.01, audioContext.currentTime + 0.5);

    oscillator.start(audioContext.currentTime);
    oscillator.stop(audioContext.currentTime + 0.5);
  } catch (error) {
    console.error('Failed to play beep sound:', error);
  }
};