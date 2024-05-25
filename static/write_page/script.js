document.getElementById('recordButton').addEventListener('click', async () => {
    const recordButton = document.getElementById('recordButton');
    const statusElement = document.getElementById('status');
    let countdown = 3;

    // Countdown logic before recording
    const countdownInterval = setInterval(() => {
        statusElement.textContent = `Recording starts in ${countdown}...`;
        countdown--;

        if (countdown < 0) {
            clearInterval(countdownInterval);
            startRecording();
            recordButton.style.display = 'none'; // Скрыть кнопку после начала записи
        }
    }, 1000);
});

async function startRecording() {
    const statusElement = document.getElementById('status');
    statusElement.textContent = 'Recording...';

    try {
        const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
        const mediaRecorder = new MediaRecorder(stream);
        const chunks = [];
        
        mediaRecorder.ondataavailable = event => chunks.push(event.data);
        
        mediaRecorder.start();

        // Recording countdown logic
        let recordingTime = 3;
        const recordingInterval = setInterval(() => {
            statusElement.textContent = `Recording... ${recordingTime}`;
            recordingTime--;

            if (recordingTime < 0) {
                clearInterval(recordingInterval);
                mediaRecorder.stop();
            }
        }, 1000);

        mediaRecorder.onstop = async () => {
            statusElement.textContent = 'Recording stopped, sending data...';
            const formData = new FormData();

            for (let i = 0; i < 3; i++) {
                const blob = new Blob(chunks, { type: 'audio/wav' });
                formData.append(`audio${i + 1}`, blob, `recording${i + 1}.wav`);
            }

            try {
                const userId = 'your-user-id'; // Replace with actual user ID
                const response = await fetch(`/api/voice/${userId}`, {
                    method: 'POST',
                    body: formData,
                });

                const result = await response.json();
                statusElement.textContent = result.message || 'Recording sent successfully!';
            } catch (error) {
                statusElement.textContent = 'Error sending recording: ' + error.message;
            }
        };
    } catch (error) {
        statusElement.textContent = 'Error accessing microphone: ' + error.message;
    }
}
